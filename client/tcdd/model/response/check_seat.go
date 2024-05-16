package response

type CheckSeatResponse struct {
	SeatLocked   bool                                `json:"koltukLocked"`
	ResponseInfo StationEmptyPlaceSearchResponseInfo `json:"cevapBilgileri"`
}
