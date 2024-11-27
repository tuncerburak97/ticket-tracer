package entity

import (
	"gorm.io/gorm"
	"time"
)

type TicketRequest struct {
	ID                  string `gorm:"primaryKey"`
	CreatedAt           time.Time
	UpdatedAt           time.Time
	DepartureDate       string `gorm:"column:departure_date"`
	DepartureStationID  int64  `gorm:"column:departure_station_id"`
	DepartureStation    string `gorm:"column:departure_station"`
	ArrivalDate         string `gorm:"column:arrival_date"`
	ArrivalStationID    int64  `gorm:"column:arrival_station_id"`
	ArrivalStation      string `gorm:"column:arrival_station"`
	TourID              int64  `gorm:"column:tour_id"`
	TrainID             int64  `gorm:"column:train_id"`
	Email               string `gorm:"column:email"`
	IsEmailNotification bool   `gorm:"column:is_email_notification"`
	Status              string `gorm:"column:status"`
	TotalAttempt        int    `gorm:"column:total_attempt"`
}

func (TicketRequest) TableName() string {
	return "ticket_request"
}

func (entity *TicketRequest) BeforeCreate(tx *gorm.DB) (err error) {
	now := time.Now()
	entity.CreatedAt = now
	entity.UpdatedAt = now
	return nil
}
