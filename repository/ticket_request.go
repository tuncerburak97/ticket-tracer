package repository

import (
	"gorm.io/gorm"
	"ticker-tracer/config/db"
	"ticker-tracer/model/entity"
)

type TicketRequestRepository interface {
	Create(ticketRequest *entity.TicketRequest) error
	FindAll() ([]entity.TicketRequest, error)
	FindById(id string) (*entity.TicketRequest, error)
	FindByStatus(status string) ([]entity.TicketRequest, error)
	Update(ticketRequest *entity.TicketRequest) error
}

type ticketRequestRepository struct {
	db *gorm.DB
}

var ticketRequestRepositoryInstance TicketRequestRepository

func NewTicketRequestRepository() TicketRequestRepository {
	return &ticketRequestRepository{db.GetDb()}
}

func GetTicketRequestRepository() TicketRequestRepository {

	if ticketRequestRepositoryInstance == nil {
		ticketRequestRepositoryInstance = NewTicketRequestRepository()
	}
	return ticketRequestRepositoryInstance
}

func (r *ticketRequestRepository) Create(ticketRequest *entity.TicketRequest) error {
	return r.db.Create(ticketRequest).Error
}
func (r *ticketRequestRepository) FindAll() ([]entity.TicketRequest, error) {
	var ticketRequests []entity.TicketRequest
	if err := r.db.Find(&ticketRequests).Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}

func (r *ticketRequestRepository) FindById(id string) (ticketRequest *entity.TicketRequest, err error) {
	var ticketRequestResponse entity.TicketRequest

	if err := r.db.Take(&ticketRequestResponse, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ticketRequestResponse, nil
}

func (r *ticketRequestRepository) FindByStatus(status string) ([]entity.TicketRequest, error) {
	var ticketRequests []entity.TicketRequest
	if err := r.db.Find(&ticketRequests, "status = ?", status).Error; err != nil {
		return nil, err
	}
	return ticketRequests, nil
}

func (r *ticketRequestRepository) Update(ticketRequest *entity.TicketRequest) error {
	return r.db.Save(ticketRequest).Error
}
