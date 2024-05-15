package train

import (
	"fmt"
	"log"
	"strconv"
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

	var emails = make([]string, 0)

	if len(ts.requests) == 0 {
		log.Printf("No request to process")
	}

	for _, searchTrainRequest := range ts.requests {
		foundedMail := ts.processRequest(searchTrainRequest)
		if foundedMail != "" {
			emails = append(emails, foundedMail)
		}
	}

	ts.RemoveRequestByEmail(emails)
}

func (ts *TrainScheduler) processRequest(request tcddServiceModel.SearchTrainRequestDetail) (email string) {

	log.Printf("Processing request: %s", request.Email)

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
	tourId, err := strconv.ParseInt(request.TourID, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	remainingDisabledNumber, found := ts.findTrip(search, tourId)
	if found {
		return ts.handleFoundTrip(request, int(remainingDisabledNumber))
	}
	log.Printf("Trip not found for request: %s", request.Email)
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

func (ts *TrainScheduler) handleFoundTrip(request tcddServiceModel.SearchTrainRequestDetail, remainingDisabledNumber int) (email string) {

	tourId, err := strconv.ParseInt(request.TourID, 10, 64)
	if err != nil {
		fmt.Println(err)
		return
	}
	placeSearch, err := ts.tcddClient.StationEmptyPlaceSearch(tcddClientRequest.StationEmptyPlaceSearchRequest{
		ChannelCode:   "3",
		Language:      0,
		TourTitleID:   tourId,
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
	if availablePlace > 0 {
		log.Printf("Found trip for request: %s, Date: %s, From: %s, To: %s",
			request.Email,
			request.DepartureDate,
			externalInfo.DepartureStation,
			externalInfo.ArrivalStation)

		log.Printf("For Request: %s Total empty place: %d and total disabled number: %d", email, totalEmptyPlace, remainingDisabledNumber)
		sendEmail(request.Email, availablePlace)
		return request.Email
	}

	log.Printf("No available place for request: %s, Date: %s, From: %s, To: %s",
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

func sendEmail(recipient string, availablePlace int) {
	{

		email := emailModel.Email{
			To:      recipient,
			Subject: "Tren Bilet Uyarısı",
			Body:    "Aradığınız trenin biletleri bulundu. Toplam boş yer sayısı:" + fmt.Sprint(availablePlace) + ". Maili aldıktan sonra tekrar bilgilendirme almak için yeni bir talepte bulunmanız gerekmektedir.",
		}

		// Send the email
		err := trainSchedulerInstance.mailClient.SendEmail(email)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}
		log.Printf("Email sent to: %s", recipient)
	}
}

func (ts *TrainScheduler) RemoveRequestByEmail(emails []string) {
	newRequests := make([]tcddServiceModel.SearchTrainRequestDetail, 0)

	for _, request := range ts.requests {
		found := false
		for _, email := range emails {
			if request.Email == email {
				found = true
				log.Printf("Removing request: %s", request.Email)
				break
			}
		}
		if !found {
			newRequests = append(newRequests, request)
		}
	}

	ts.requests = newRequests
}
