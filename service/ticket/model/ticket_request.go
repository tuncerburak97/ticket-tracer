package model

type RetrieveTicketRequest struct {
	ID               string `json:"id"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	DepartureStation string `json:"departure_station"`
	DepartureDate    string `json:"departure_date"`
	ArrivalStation   string `json:"arrival_station"`
	ArrivalDate      string `json:"arrival_date"`
	Mail             string `json:"mail"`
	Status           string `json:"status"`
	TotalAttempt     int    `json:"total_attempt"`
}
