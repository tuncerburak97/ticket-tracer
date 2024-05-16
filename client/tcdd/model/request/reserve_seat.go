package request

type ReserveSeatRequest struct {
	ChannelCode        string `json:"kanalKodu"`
	Language           int    `json:"dil"`
	TourID             int    `json:"seferId"`
	WagonOrder         int    `json:"vagonSiraNo"`
	SeatNo             string `json:"koltukNo"`
	Gender             string `json:"cinsiyet"`
	ArrivalStationID   int    `json:"inisIst"`
	DepartureStationID int    `json:"binisIst"`
	Minute             int    `json:"dakika"`
	Huawei             bool   `json:"huawei"`
}
