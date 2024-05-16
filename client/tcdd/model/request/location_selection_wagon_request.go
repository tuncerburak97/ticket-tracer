package request

type LocationSelectionWagonRequest struct {
	ChannelCode          string `json:"kanalKodu"`
	Language             int    `json:"dil"`
	TourTitleID          string `json:"seferBaslikId"`
	WagonOrderNo         int    `json:"vagonSiraNo"`
	DepartureStationName string `json:"binisIst"`
	ArrivalStationName   string `json:"InisIst"`
}
