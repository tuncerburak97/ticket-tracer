package tcdd

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"sync"
	"ticker-tracer/client/tcdd"
	clientRequestModel "ticker-tracer/client/tcdd/model/request"
	clientResponseModel "ticker-tracer/client/tcdd/model/response"
	"ticker-tracer/scheduler/train"
	serviceModel "ticker-tracer/service/tcdd/model"
	"time"
)

type TccdService struct {
	tcddClient     *tcdd.TcddHttpClient
	trainScheduler *train.TrainScheduler
	stations       *clientResponseModel.StationLoadResponse
	once           sync.Once
}

type TccdServiceInterface interface {
	GetStations() (*clientResponseModel.StationLoadResponse, error)
	LoadStations() (*serviceModel.StationInformation, error)
	AddSearchRequest(request *serviceModel.SearchTrainRequest) (*serviceModel.SearchTrainResponse, error)
	QueryTrain(request *serviceModel.QueryTrainRequest) (*serviceModel.QueryTrainResponse, error)
}

func NewTcddService() *TccdService {
	return &TccdService{
		tcddClient:     tcdd.GetTcddHttpClientInstance(),
		trainScheduler: train.GetTrainSchedulerInstance(),
	}
}

func (ts *TccdService) GetStations() (*clientResponseModel.StationLoadResponse, error) {
	var err error
	ts.once.Do(func() {
		stationLoadRequest := clientRequestModel.StationLoadRequest{
			Language:    0,
			ChannelCode: "3",
			Date:        "Nov 10, 2011 12:00:00 AM",
			SalesQuery:  true,
		}
		ts.stations, err = ts.tcddClient.LoadAllStation(stationLoadRequest)
	})
	return ts.stations, err
}

func (ts *TccdService) LoadStations() (*serviceModel.StationInformation, error) {
	stations, err := ts.GetStations()
	if err != nil {
		return &serviceModel.StationInformation{
			Message:  "Error loading stations",
			Success:  false,
			Response: make([]serviceModel.LoadStationResponse, 0),
		}, err
	}
	var stationList []serviceModel.LoadStationResponse

	for _, station := range stations.StationInformation {
		isYht := false
		for _, stationTrainType := range station.StationTrainTypes {
			if stationTrainType == "YHT" {
				isYht = true
			}
		}
		if !isYht {
			continue
		}

		var toStationList []serviceModel.ToStationList
		for _, toStation := range station.ToStationIDs {
			toStationData, _ := GetStationByStationID(stations.StationInformation, toStation)
			if toStationData == nil {
				continue
			}
			toStationList = append(toStationList, serviceModel.ToStationList{
				ToStationID:   toStation,
				ToStationName: toStationData.StationName,
			})
			sortToStationListByName(toStationList)
		}
		stationList = append(stationList, serviceModel.LoadStationResponse{
			StationID:         station.StationID,
			StationCode:       station.StationCode,
			StationName:       station.StationName,
			StationViewName:   station.StationViewName,
			StationTrainTypes: station.StationTrainTypes,
			ToStationList:     toStationList,
		})

	}

	sortStationsByStationName(stationList)

	return &serviceModel.StationInformation{
		Message:  "Stations loaded",
		Success:  true,
		Response: stationList,
	}, nil
}
func sortStationsByStationName(loadStationResponse []serviceModel.LoadStationResponse) {
	sort.Slice(loadStationResponse, func(i, j int) bool {
		return loadStationResponse[i].StationName < loadStationResponse[j].StationName
	})
}

func sortToStationListByName(toStationList []serviceModel.ToStationList) {
	sort.Slice(toStationList, func(i, j int) bool {
		return toStationList[i].ToStationName < toStationList[j].ToStationName
	})

}

func (ts *TccdService) AddSearchRequest(requests *serviceModel.SearchTrainRequest) (*serviceModel.SearchTrainResponse, error) {
	for _, request := range requests.Request {
		parsedTime, err := time.Parse("Jan 2, 2006 03:04:05 PM", request.DepartureDate)
		if err != nil {
			return nil, fmt.Errorf("invalid departure date: %v", err)
		}
		var now = time.Now()
		if now.After(parsedTime) {
			return nil, errors.New("past departure date")
		}

		if !validateEmail(request.Email) {
			return nil, errors.New("invalid email format")
		}
		if stations, err := ts.GetStations(); err != nil {
			return nil, fmt.Errorf("error getting stations: %v", err)
		} else {

			if !checkStationIDIsValid(request.DepartureStationID, stations.StationInformation) || !checkStationIDIsValid(request.ArrivalStationID, stations.StationInformation) {
				return nil, errors.New("invalid arrival or departure station id")
			}

			departureStation, err := GetStationByStationID(stations.StationInformation, request.DepartureStationID)
			if err != nil {
				return nil, fmt.Errorf("error getting departure station: %v", err)
			}
			found := false
			for _, toStationID := range departureStation.ToStationIDs {
				if toStationID == request.ArrivalStationID {
					found = true
				}
			}
			if !found {
				return nil, errors.New("arrival station is not reachable from departure station")
			}

			arrivalStation, _ := GetStationByStationID(stations.StationInformation, request.ArrivalStationID)

			externalInfo := serviceModel.ExternalInformation{
				DepartureStation: departureStation.StationName,
				ArrivalStation:   arrivalStation.StationName,
				DepartureDate:    request.DepartureDate,
			}
			newRequest := serviceModel.SearchTrainRequestDetail{
				DepartureDate:       request.DepartureDate,
				DepartureStationID:  request.DepartureStationID,
				ArrivalStationID:    request.ArrivalStationID,
				TourID:              request.TourID,
				TrainID:             request.TrainID,
				Email:               request.Email,
				IsEmailNotification: request.IsEmailNotification,
				ExternalInformation: externalInfo,
			}
			err = checkEmailRequestExceedThreshold(request.Email, *requests)
			if err != nil {
				return nil, err
			}
			ts.trainScheduler.AddRequest(newRequest)

		}
	}
	return &serviceModel.SearchTrainResponse{
		Message: "Request added to scheduler",
		Success: true,
	}, nil
}

func (ts *TccdService) QueryTrain(request *serviceModel.QueryTrainRequest) (*serviceModel.QueryTrainResponse, error) {
	criteria := clientRequestModel.Criteria{
		SalesChannel:       3,
		DepartureStation:   request.DepartureStationName,
		IsMapDeparture:     false,
		ArrivalStation:     request.ArrivalStationName,
		IsMapArrival:       false,
		DepartureDate:      request.DepartureDate,
		IsRegional:         false,
		OperationType:      0,
		PassengerCount:     1,
		IsTransfer:         true,
		DepartureStationID: request.DepartureStationID,
		ArrivalStationID:   request.ArrivalStationID,
		TravelType:         1,
	}

	tripSearchResponse, err := ts.tcddClient.TripSearch(clientRequestModel.TripSearchRequest{
		ChannelCode: 3,
		Language:    0,
		Criteria:    criteria,
	})
	if err != nil {
		return nil, fmt.Errorf("error querying train: %v", err)
	}

	var wg sync.WaitGroup
	detailsChan := make(chan serviceModel.QueryTrainResponseDetail)
	errChan := make(chan error, len(tripSearchResponse.SearchResult))

	for _, trip := range tripSearchResponse.SearchResult {
		wg.Add(1)
		go ts.processTripSearchResult(&wg, detailsChan, errChan, trip, request, tripSearchResponse)
	}

	go func() {
		wg.Wait()
		close(detailsChan)
		close(errChan)
	}()

	var details []serviceModel.QueryTrainResponseDetail
	for detail := range detailsChan {
		details = append(details, detail)
	}

	if len(errChan) > 0 {
		return nil, <-errChan
	}

	orderByArrivalDate(details)

	return &serviceModel.QueryTrainResponse{
		Details: details,
	}, nil
}

func orderByArrivalDate(Details []serviceModel.QueryTrainResponseDetail) {
	sort.Slice(Details, func(i, j int) bool {
		iTime, _ := time.Parse("Jan 2, 2006 03:04:05 PM", Details[i].ArrivalDate)
		jTime, _ := time.Parse("Jan 2, 2006 03:04:05 PM", Details[j].ArrivalDate)
		return iTime.Before(jTime)
	})
}

func (ts *TccdService) processTripSearchResult(
	wg *sync.WaitGroup,
	detailsChan chan<- serviceModel.QueryTrainResponseDetail,
	errChan chan<- error,
	trip clientResponseModel.SearchResult,
	request *serviceModel.QueryTrainRequest,
	tripSearchResponse *clientResponseModel.TripSearchResponse,
) {
	defer wg.Done()

	remainingDisabledNumber, _ := findTrip(tripSearchResponse, trip.TourID)
	placeSearch, err := ts.tcddClient.StationEmptyPlaceSearch(clientRequestModel.StationEmptyPlaceSearchRequest{
		ChannelCode:   "3",
		Language:      0,
		TourTitleID:   trip.TourID,
		DepartureStID: request.DepartureStationID,
		ArrivalStID:   int(request.ArrivalStationID),
	})
	if err != nil {
		errChan <- fmt.Errorf("error getting empty place: %v", err)
		return
	}

	totalEmptyPlace := calculateTotalEmptyPlace(placeSearch.EmptyPlaceList)
	detailsChan <- serviceModel.QueryTrainResponseDetail{
		TrainID:            trip.TrainID,
		TrainName:          trip.TrainName,
		TrainCode:          trip.TrainCode,
		TourID:             trip.TourID,
		DepartureDate:      trip.DepartureDate,
		ArrivalDate:        trip.ArrivalDate,
		ArrivalStation:     trip.ArrivalStation,
		DepartureStation:   trip.DepartureStation,
		ArrivalStationID:   trip.ArrivalStationID,
		DepartureStationID: trip.DepartureStationID,
		EmptyPlace: serviceModel.EmptyPlace{
			DisabledPlaceCount:          remainingDisabledNumber,
			TotalEmptyPlaceCount:        int64(totalEmptyPlace),
			NormalPeopleEmptyPlaceCount: int64(totalEmptyPlace) - remainingDisabledNumber,
		},
	}
}

func calculateTotalEmptyPlace(emptyPlaceList []clientResponseModel.EmptyPlace) int {
	totalEmptyPlace := 0
	for _, emptyPlace := range emptyPlaceList {
		totalEmptyPlace += emptyPlace.EmptyPlace
	}
	return totalEmptyPlace
}
func findTrip(search *clientResponseModel.TripSearchResponse, tourID int64) (int64, bool) {
	for _, trip := range search.SearchResult {
		if trip.TourID == tourID {
			if len(trip.WagonTypesEmptyPlace) > 0 {
				return trip.WagonTypesEmptyPlace[0].RemainingDisabledNumber, true
			}
		}
	}
	return 0, false
}

// commons
func validateEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
	validationResult := regexp.MustCompile(emailRegex).MatchString(email)
	return validationResult
}
func GetStationByStationID(stations []clientResponseModel.StationInformation, stationID int64) (*clientResponseModel.StationInformation, error) {
	for _, station := range stations {
		if station.StationID == stationID {
			return &station, nil
		}
	}
	return nil, fmt.Errorf("no station found with ID: %v", stationID)
}

func checkStationIDIsValid(stationID int64, stations []clientResponseModel.StationInformation) bool {
	for _, station := range stations {
		if station.StationID == stationID {
			return true
		}
	}
	return false

}

func checkEmailRequestExceedThreshold(email string, requests serviceModel.SearchTrainRequest) error {
	foundedCount := 0
	for _, request := range requests.Request {
		if request.Email == email {
			foundedCount++
		}
	}
	if foundedCount > 5 {
		return errors.New("exceed threshold")
	}
	return nil
}
