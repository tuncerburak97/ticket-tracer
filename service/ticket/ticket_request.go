package ticket

import (
	"ticker-tracer/repository"
	"ticker-tracer/service/ticket/model"
)

type TicketRequestService struct {
	ticketRequestRepository repository.TicketRequestRepository
}

type TicketRequestServiceInterface interface {
	FindById(id string) (model.RetrieveTicketRequest, error)
	FindAll() ([]model.RetrieveTicketRequest, error)
	FindByMail(mail string) ([]model.RetrieveTicketRequest, error)
	FindByStatus(status string) ([]model.RetrieveTicketRequest, error)
}

var ticketRequestService *TicketRequestService

func NewTicketRequestService() *TicketRequestService {
	ticketRequestService = &TicketRequestService{
		ticketRequestRepository: repository.GetTicketRequestRepository(),
	}
	return ticketRequestService
}

func GetTicketRequestService() *TicketRequestService {
	if ticketRequestService == nil {
		return NewTicketRequestService()
	}
	return ticketRequestService
}

func (service *TicketRequestService) FindById(id string) (model.RetrieveTicketRequest, error) {
	entity, err := service.ticketRequestRepository.FindById(id)
	if err != nil {
		return model.RetrieveTicketRequest{}, err
	}

	dto := model.RetrieveTicketRequest{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
	}
}
