package request

type StationEmptyPlaceSearchRequest struct {
	ChannelCode   string `json:"kanalKodu"`
	Language      int    `json:"dil"`
	TourTitleID   int64  `json:"seferBaslikId"`
	DepartureStID int64  `json:"binisIstId"`
	ArrivalStID   int    `json:"inisIstId"`
}
