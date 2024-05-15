package request

type StationLoadRequest struct {
	ChannelCode string `json:"kanalKodu"`
	Language    int    `json:"dil"`
	Date        string `json:"tarih"`
	SalesQuery  bool   `json:"satisSorgu"`
}
