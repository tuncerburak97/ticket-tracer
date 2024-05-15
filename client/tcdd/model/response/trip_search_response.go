package response

type TripSearchResponse struct {
	TripSearchResponseInfo TripSearchResponseInfo `json:"cevapBilgileri"`
	SearchResult           []SearchResult         `json:"seferSorgulamaSonucList"`
}

type TripSearchResponseInfo struct {
	ResponseCode string      `json:"cevapKodu"`
	ResponseMsg  string      `json:"cevapMsj"`
	Detail       interface{} `json:"detay"`
}

type SearchResult struct {
	TrainID              int64                  `json:"trenId"`
	TrainName            string                 `json:"trenAdi"`
	TrainType            string                 `json:"trenTipi"`
	TrainCode            string                 `json:"trenKodu"`
	TrainTourTktID       int64                  `json:"trenTuruTktId"`
	TourName             string                 `json:"seferAdi"`
	TourID               int64                  `json:"seferId"`
	DepartureDate        string                 `json:"binisTarih"`
	ArrivalDate          string                 `json:"inisTarih"`
	WagonTypesEmptyPlace []WagonTypesEmptyPlace `json:"vagonTipleriBosYerUcret"`
	SalesStatus          int64                  `json:"satisDurum"`
	DepartureStationID   int64                  `json:"binisIstasyonId"`
	ArrivalStationID     int64                  `json:"inisIstasyonId"`
	ArrivalStation       string                 `json:"inisIstasyonu"`
	DepartureStation     string                 `json:"binisIstasyonu"`
}

type WagonTypesEmptyPlace struct {
	WagonType               string  `json:"vagonTip"`
	RemainingNumber         int64   `json:"kalanSayi"`
	CheapTicketPrice        float64 `json:"hesapliBiletFiyati"`
	StandardTicketPrice     float64 `json:"standartBiletFiyati"`
	RemainingDisabledNumber int64   `json:"kalanEngelliKoltukSayisi"`
}
