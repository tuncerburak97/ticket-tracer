package ticket

import (
	"ticker-tracer/repository"
	"ticker-tracer/service/ticket/model"
)

type RequestService struct {
	ticketRequestRepository repository.TicketRequestRepository
}

type ServiceInterface interface {
	FindById(id string) (model.RetrieveTicketRequest, error)
	FindAll() ([]model.RetrieveTicketRequest, error)
	FindByMail(mail string) ([]model.RetrieveTicketRequest, error)
	FindByStatus(status string) ([]model.RetrieveTicketRequest, error)
	FindByMailAndStatus(status, mail string) ([]model.RetrieveTicketRequest, error)
}

var ticketRequestService *RequestService

func NewTicketRequestService() *RequestService {
	ticketRequestService = &RequestService{
		ticketRequestRepository: repository.GetTicketRequestRepository(),
	}
	return ticketRequestService
}

func GetTicketRequestService() *RequestService {
	if ticketRequestService == nil {
		return NewTicketRequestService()
	}
	return ticketRequestService
}

func (service *RequestService) FindById(id string) (model.RetrieveTicketRequest, error) {
	entity, err := service.ticketRequestRepository.FindById(id)
	if err != nil {
		return model.RetrieveTicketRequest{}, err
	}

	dto := model.RetrieveTicketRequest{
		ID:               entity.ID,
		CreatedAt:        entity.CreatedAt,
		UpdatedAt:        entity.UpdatedAt,
		DepartureStation: entity.DepartureStation,
		DepartureDate:    entity.DepartureDate,
		ArrivalStation:   entity.ArrivalStation,
		ArrivalDate:      entity.ArrivalDate,
		Email:            entity.Email,
		Status:           entity.Status,
		TotalAttempt:     entity.TotalAttempt,
	}

	return dto, nil
}

func (service *RequestService) FindAll() ([]model.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindAll()
	if err != nil {
		return nil, err
	}

	dtoList := []model.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := model.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *RequestService) FindByMail(mail string) ([]model.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByMail(mail)
	if err != nil {
		return nil, err
	}

	dtoList := []model.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := model.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *RequestService) FindByStatus(status string) ([]model.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByStatus(status)
	if err != nil {
		return nil, err
	}

	dtoList := []model.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := model.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}

func (service *RequestService) FindByMailAndStatus(mail, status string) ([]model.RetrieveTicketRequest, error) {
	entityList, err := service.ticketRequestRepository.FindByMailAndStatus(mail, status)
	if err != nil {
		return nil, err
	}

	dtoList := []model.RetrieveTicketRequest{}
	for _, entity := range entityList {
		dto := model.RetrieveTicketRequest{
			ID:               entity.ID,
			CreatedAt:        entity.CreatedAt,
			UpdatedAt:        entity.UpdatedAt,
			DepartureStation: entity.DepartureStation,
			DepartureDate:    entity.DepartureDate,
			ArrivalStation:   entity.ArrivalStation,
			ArrivalDate:      entity.ArrivalDate,
			Email:            entity.Email,
			Status:           entity.Status,
			TotalAttempt:     entity.TotalAttempt,
		}
		dtoList = append(dtoList, dto)
	}

	return dtoList, nil
}
