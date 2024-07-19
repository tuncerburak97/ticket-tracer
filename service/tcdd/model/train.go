package model

type QueryTrainRequest struct {
	DepartureStationID   int64  `json:"binisIstasyonId"`
	DepartureStationName string `json:"binisIstasyon"`
	ArrivalStationID     int64  `json:"inisIstasyonId"`
	ArrivalStationName   string `json:"inisIstasyonu"`
	DepartureDate        string `json:"gidisTarih"`
}

type QueryTrainResponse struct {
	Details []QueryTrainResponseDetail `json:"details"`
}
type QueryTrainResponseDetail struct {
	TrainID            int64      `json:"trainID"`
	TrainName          string     `json:"trainName"`
	TrainCode          string     `json:"trainCode"`
	TourID             int64      `json:"tourID"`
	DepartureDate      string     `json:"departureDate"`
	ArrivalDate        string     `json:"arrivalDate"`
	EmptyPlace         EmptyPlace `json:"emptyPlace"`
	ArrivalStation     string     `json:"arrivalStation"`
	DepartureStation   string     `json:"departureStation"`
	DepartureStationID int64      `json:"departureStationID"`
	ArrivalStationID   int64      `json:"arrivalStationID"`
}

type EmptyPlace struct {
	DisabledPlaceCount          int64 `json:"disabledPlaceCount"`
	TotalEmptyPlaceCount        int64 `json:"totalEmptyPlaceCount"`
	NormalPeopleEmptyPlaceCount int64 `json:"normalPeopleEmptyPlaceCount"`
}

type SearchTrainRequest struct {
	Request []SearchTrainRequestDetail `json:"request"`
}

type SearchTrainRequestDetail struct {
	RequestID           string              `json:"requestID"`
	DepartureDate       string              `json:"gidisTarih"`
	DepartureStationID  int64               `json:"binisIstasyonId"`
	ArrivalStationID    int64               `json:"inisIstasyonId"`
	ArrivalDate         string              `json:"inisTarih"`
	TourID              int64               `json:"tourID"`
	TrainID             int64               `json:"trainID"`
	Email               string              `json:"email"`
	IsEmailNotification bool                `json:"emailNotification"`
	ExternalInformation ExternalInformation `json:"externalInformation"`
}

type SearchTrainResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ExternalInformation struct {
	DepartureStation string `json:"departureStation"`
	ArrivalStation   string `json:"arrivalStation"`
	DepartureDate    string `json:"departureDate"`
	ArrivalDate      string `json:"arrivalDate"`
}

type StationInformation struct {
	Response []LoadStationResponse `json:"response"`
	Message  string                `json:"message"`
	Success  bool                  `json:"success"`
}

type LoadStationResponse struct {
	StationName       string          `json:"stationName"`
	StationID         int64           `json:"stationID"`
	StationCode       string          `json:"stationCode"`
	StationTrainTypes []string        `json:"stationTrainTypes"`
	StationViewName   string          `json:"stationViewName"`
	ToStationList     []ToStationList `json:"toStationList"`
}

type ToStationList struct {
	ToStationID   int64  `json:"toStationId"`
	ToStationName string `json:"toStationName"`
}
