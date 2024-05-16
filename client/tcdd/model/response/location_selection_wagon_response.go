package response

type LocationSelectionWagonResponse struct {
	LocationSelectionWagonResponseData LocationSelectionWagonResponseData  `json:"vagonHaritasiIcerikDVO"`
	ResponseInfo                       StationEmptyPlaceSearchResponseInfo `json:"cevapBilgileri"`
}

type LocationSelectionWagonResponseData struct {
	SeatInformationList []SeatInformation `json:"koltukDurumlari"`
}

type SeatInformation struct {
	Status       int    `json:"durum"`
	SeatNo       string `json:"koltukNo"`
	TourTitleID  int    `json:"seferBaslikId"`
	WagonOrderNo int    `json:"vagonSiraNo"`
}
