package response

type StationEmptyPlaceSearchResponse struct {
	ResponseInfo   StationEmptyPlaceSearchResponseInfo `json:"cevapBilgileri"`
	EmptyPlaceList []EmptyPlace                        `json:"vagonBosYerList"`
}

type StationEmptyPlaceSearchResponseInfo struct {
	ResponseCode string `json:"cevapKodu"`
	ResponseMsg  string `json:"cevapMsj"`
	Detail       string `json:"detay"`
}

type EmptyPlace struct {
	WagonTitleID int64 `json:"vagonBaslikId"`
	WagonOrderNo int   `json:"vagonSiraNo"`
	EmptyPlace   int   `json:"bosYer"`
}
