package train

import (
	"fmt"
	"log"
	"sync"
	"ticker-tracer/client/notification/mail"
	emailModel "ticker-tracer/client/notification/mail/model"
	"ticker-tracer/client/tcdd"
	tcddClientRequest "ticker-tracer/client/tcdd/model/request"
	tcddClientResponse "ticker-tracer/client/tcdd/model/response"
	tcddServiceModel "ticker-tracer/service/tcdd/model"
)

type TrainScheduler struct {
	tcddClient *tcdd.TcddHttpClient
	mailClient *mail.MailHttpClient
	stations   *tcddClientResponse.StationLoadResponse
	once       sync.Once
	requests   []tcddServiceModel.SearchTrainRequestDetail
	mu         sync.Mutex
}

var trainSchedulerInstance *TrainScheduler

func GetTrainSchedulerInstance() *TrainScheduler {
	if trainSchedulerInstance == nil {
		trainSchedulerInstance = NewTrainScheduler(tcdd.GetTcddHttpClientInstance(),
			mail.GetMailHttpClientInstance())
	}
	return trainSchedulerInstance

}

func NewTrainScheduler(tcddClient *tcdd.TcddHttpClient,
	mailClient *mail.MailHttpClient,
) *TrainScheduler {
	return &TrainScheduler{
		tcddClient: tcddClient,
		mailClient: mailClient,
		requests:   make([]tcddServiceModel.SearchTrainRequestDetail, 0),
	}
}

func (ts *TrainScheduler) AddRequest(request tcddServiceModel.SearchTrainRequestDetail) {
	ts.requests = append(ts.requests, request)
}

func (ts *TrainScheduler) getStations() (*tcddClientResponse.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		stationLoadRequest := tcddClientRequest.StationLoadRequest{
			Language:    0,
			ChannelCode: "3",
			Date:        "Nov 10, 2011 12:00:00 AM",
			SalesQuery:  true,
		}
		ts.stations, err = ts.tcddClient.LoadAllStation(stationLoadRequest)
	})
	return ts.stations, err
}

func (ts *TrainScheduler) Run() {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	log.Printf("Running train scheduler with %d requests", len(ts.requests))

	if _, err := ts.getStations(); err != nil {
		log.Printf("Error getting stations: %v", err)
		return
	}

	var foundedRequestIDList = make([]string, 0)

	if len(ts.requests) == 0 {
		log.Printf("No request to process")
	}

	for _, searchTrainRequest := range ts.requests {
		foundedRequestIDS := ts.processRequest(searchTrainRequest)
		if foundedRequestIDS != "" {
			foundedRequestIDList = append(foundedRequestIDList, foundedRequestIDS)
		}
	}

	ts.RemoveFoundedRequestByRequestID(foundedRequestIDList)
}

func (ts *TrainScheduler) processRequest(request tcddServiceModel.SearchTrainRequestDetail) (requestID string) {

	log.Printf("Processing request: %s", request.RequestID)

	criteria := tcddClientRequest.Criteria{
		SalesChannel:       3,
		DepartureStation:   request.ExternalInformation.DepartureStation,
		IsMapDeparture:     false,
		ArrivalStation:     request.ExternalInformation.ArrivalStation,
		IsMapArrival:       false,
		DepartureDate:      request.ExternalInformation.DepartureDate,
		IsRegional:         false,
		OperationType:      0,
		PassengerCount:     1,
		IsTransfer:         true,
		DepartureStationID: request.DepartureStationID,
		ArrivalStationID:   request.ArrivalStationID,
		TravelType:         1,
	}

	search, err := ts.tcddClient.TripSearch(tcddClientRequest.TripSearchRequest{
		ChannelCode: 3,
		Language:    0,
		Criteria:    criteria,
	})
	if err != nil {
		log.Printf("Error searching trip: %v", err)
		return
	}

	b := search.TripSearchResponseInfo.ResponseCode != "000"
	if b {
		log.Printf("Error searching trip: %v", search.TripSearchResponseInfo.ResponseMsg)
		return
	}
	remainingDisabledNumber, found := ts.findTrip(search, request.TourID)
	if found {
		return ts.handleFoundTrip(request, int(remainingDisabledNumber), search.SearchResult[0].ArrivalDate)
	}
	log.Printf("Trip not found for request: %s and email: %s date: %s from: %s to: %s",
		request.RequestID,
		request.Email,
		request.DepartureDate,
		request.ExternalInformation.DepartureStation,
		request.ExternalInformation.ArrivalStation)

	return ""

}

func (ts *TrainScheduler) findTrip(search *tcddClientResponse.TripSearchResponse, tourID int64) (int64, bool) {
	for _, trip := range search.SearchResult {
		if trip.TourID == tourID {
			if len(trip.WagonTypesEmptyPlace) > 0 {
				return trip.WagonTypesEmptyPlace[0].RemainingDisabledNumber, true
			}
		}
	}
	return 0, false
}

func (ts *TrainScheduler) handleFoundTrip(request tcddServiceModel.SearchTrainRequestDetail, remainingDisabledNumber int, arrivalDate string) (requestID string) {

	placeSearch, err := ts.tcddClient.StationEmptyPlaceSearch(tcddClientRequest.StationEmptyPlaceSearchRequest{
		ChannelCode:   "3",
		Language:      0,
		TourTitleID:   request.TourID,
		DepartureStID: request.DepartureStationID,
		ArrivalStID:   int(request.ArrivalStationID),
	})
	if err != nil {
		log.Printf("Error getting empty place: %v", err)
		return
	}

	totalEmptyPlace := calculateTotalEmptyPlace(placeSearch.EmptyPlaceList)
	availablePlace := totalEmptyPlace - remainingDisabledNumber
	externalInfo := request.ExternalInformation
	externalInfo.ArrivalDate = arrivalDate
	if availablePlace > 0 {
		log.Printf("Found trip for request: %s,Email: %s Date: %s, From: %s, To: %s",
			request.RequestID,
			request.Email,
			request.DepartureDate,
			externalInfo.DepartureStation,
			externalInfo.ArrivalStation)

		log.Printf("For Request: %s with Email: %s, Date: %s, From: %s, To: %s, Total Empty Place: %d, Remaining Disabled Number: %d",
			requestID,
			request.Email,
			request.ExternalInformation.DepartureDate,
			request.ExternalInformation.DepartureStation,
			request.ExternalInformation.ArrivalStation,
			totalEmptyPlace,
			remainingDisabledNumber)
		sendEmail(
			request.Email,
			availablePlace,
			externalInfo.DepartureDate,
			externalInfo.ArrivalDate,
			externalInfo.DepartureStation, externalInfo.ArrivalStation)
		return request.RequestID
	}

	log.Printf("No available place for request: %s, Email: %s, Date: %s, From: %s, To: %s",
		request.RequestID,
		request.Email,
		request.DepartureDate,
		externalInfo.DepartureStation,
		externalInfo.ArrivalStation)

	return ""
}

func calculateTotalEmptyPlace(emptyPlaceList []tcddClientResponse.EmptyPlace) int {
	totalEmptyPlace := 0
	for _, emptyPlace := range emptyPlaceList {
		totalEmptyPlace += emptyPlace.EmptyPlace
	}
	return totalEmptyPlace
}

func sendEmail(recipient string,
	availablePlace int,
	departureDate string,
	arrivalDate string,
	departureStation string,
	arrivalStation string,
) {

	{
		body := fmt.Sprintf(`
  <html>
  <head>
  <style>
  table {
    font-family: Arial, sans-serif;
    border-collapse: collapse;
    width: 100%%;
  }

  td, th {
    border: 1px solid #dddddd;
    text-align: left;
    padding: 8px;
  }

  tr:nth-child(even) {
    background-color: #dddddd;
  }
  </style>
  </head>
  <body>
  <p>Merhaba,</p>
  <p>Aradığınız trende boş yer bulundu. &#128522;</p>
  <table>
    <tr>
      <th>Kalan Boş Yer Sayısı</th>
      <th>Kalkış Zamanı</th>
      <th>Varış Zamanı</th>
      <th>Kalkış İstasyonu</th>
      <th>Varış İstasyonu</th>
    </tr>
    <tr>
      <td>%d</td>
      <td>%s</td>
	  <td>%s</td>
      <td>%s</td>
      <td>%s</td>
    </tr>
  </table>
  <p>Tekrardan bu yolculuga dair bildirimleri takip etmek isterseniz uygulama üzerinden aynı talebi oluşturabilirsiniz</p>
  <p>İyi yolculuklar dileriz!</p>
  </body>
  </html>`, availablePlace, departureDate, arrivalDate, departureStation, arrivalStation)

		email := emailModel.Email{
			To:      recipient,
			Subject: "Tren Bilet Uyarısı",
			Body:    body,
		}

		// Send the email
		err := trainSchedulerInstance.mailClient.SendEmail(email)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}
		log.Printf("Email sent to: %s", recipient)
	}
}

func (ts *TrainScheduler) RemoveFoundedRequestByRequestID(foundedRequestIDList []string) {
	newRequests := make([]tcddServiceModel.SearchTrainRequestDetail, 0)

	for _, request := range ts.requests {
		found := false
		for _, foundedRequestID := range foundedRequestIDList {
			if request.RequestID == foundedRequestID {
				found = true
				log.Printf("Removing request: %s with Email: %s, Date: %s, From: %s, To: %s",
					request.RequestID,
					request.Email,
					request.DepartureDate,
					request.ExternalInformation.DepartureStation,
					request.ExternalInformation.ArrivalStation)
				break
			}
		}
		if !found {
			newRequests = append(newRequests, request)
		}
	}

	ts.requests = newRequests
}
