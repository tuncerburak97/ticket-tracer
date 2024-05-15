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
	log.Printf("Running train scheduler")
	if _, err := ts.getStations(); err != nil {
		log.Printf("Error getting stations: %v", err)
		return
	}

	for _, searchTrainRequest := range ts.requests {
		ts.processRequest(searchTrainRequest)
	}
}

func (ts *TrainScheduler) processRequest(request tcddServiceModel.SearchTrainRequestDetail) {

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

	remainingDisabledNumber, found := ts.findTrip(search, request.TourID)
	if found {
		ts.handleFoundTrip(request, int(remainingDisabledNumber))
	}
	if !found {
		log.Printf("Trip not found for request: %s", request.Email)
	}
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

func (ts *TrainScheduler) handleFoundTrip(request tcddServiceModel.SearchTrainRequestDetail, remainingDisabledNumber int) {
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
	log.Printf("Found trip for request: %s, Date: %s, From: %d, To: %d", request.Email, request.DepartureDate, request.DepartureStationID, request.ArrivalStationID)
	log.Printf("Total empty place: %d and total disabled number: %d", totalEmptyPlace, remainingDisabledNumber)

	availablePlace := totalEmptyPlace - remainingDisabledNumber
	if availablePlace > -1 {
		sendEmail(request.Email, availablePlace)
	}
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
			Body:    "Aradığınız trenin biletleri bulundu. Toplam boş yer sayısı:" + fmt.Sprint(availablePlace) + ". Acele edin!",
		}

		// Send the email
		err := trainSchedulerInstance.mailClient.SendEmail(email)
		if err != nil {
			fmt.Println("Error sending email:", err)
		}
	}
}
