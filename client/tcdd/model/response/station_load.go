package response

type StationLoadResponse struct {
	ResponseInfo       StationLoadResponseInfo `json:"cevapBilgileri"`
	StationInformation []StationInformation    `json:"istasyonBilgileriList"`
}

type StationInformation struct {
	StationID         int64    `json:"istasyonId"`
	StationCode       string   `json:"istasyonKodu"`
	StationName       string   `json:"istasyonAdi"`
	StationStatus     bool     `json:"istasyonDurumu"`
	Date              string   `json:"tarih"`
	ToStationIDs      []int64  `json:"toStationIds"`
	RegionCode        string   `json:"bolgeKodu"`
	StationTrainTypes []string `json:"stationTrainTypes"`
	StationViewName   string   `json:"stationViewName"`
}
type StationLoadResponseInfo struct {
	ResponseCode string `json:"cevapKodu"`
	ResponseMsg  string `json:"cevapMsj"`
}
