package request

type TripSearchRequest struct {
	ChannelCode int      `json:"kanalKodu"`
	Language    int      `json:"dil"`
	Criteria    Criteria `json:"seferSorgulamaKriterWSDVO"`
}

type Criteria struct {
	SalesChannel       int64  `json:"satisKanali"`
	DepartureStation   string `json:"binisIstasyonu"`
	IsMapDeparture     bool   `json:"binisIstasyonu_isHaritaGosterimi"`
	ArrivalStation     string `json:"inisIstasyonu"`
	IsMapArrival       bool   `json:"inisIstasyonu_isHaritaGosterimi"`
	TravelType         int64  `json:"seyahatTuru"`
	DepartureDate      string `json:"gidisTarih"`
	IsRegional         bool   `json:"bolgeselGelsin"`
	OperationType      int64  `json:"islemTipi"`
	PassengerCount     int64  `json:"yolcuSayisi"`
	IsTransfer         bool   `json:"aktarmalarGelsin"`
	DepartureStationID int64  `json:"binisIstasyonId"`
	ArrivalStationID   int64  `json:"inisIstasyonId"`
}
