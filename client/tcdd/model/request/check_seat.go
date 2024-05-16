package request

type CheckSeatRequest struct {
	TourId                  string `json:"seferId"`
	ChannelCode             string `json:"kanalKodu"`
	SelectedSeatWagonNumber int    `json:"seciliVagonSiraNo"`
	SelectedSeatNumber      string `json:"koltukNo"`
	Language                int    `json:"dil"`
}
