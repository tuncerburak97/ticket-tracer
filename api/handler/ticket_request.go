package handler

import (
	"github.com/gofiber/fiber/v2"
	"ticker-tracer/service/ticket"
	utils "ticker-tracer/util/response"
)

type TicketRequestHandlerInterface interface {
	FindRequestById(c *fiber.Ctx) error
	FindAllRequest(c *fiber.Ctx) error
	FindRequestByMail(c *fiber.Ctx) error
	FindRequestByStatus(c *fiber.Ctx) error
	FindRequestByMailAndStatus(c *fiber.Ctx) error
}

type TicketRequestHandler struct {
	s *ticket.RequestService
}

func NewTicketRequestHandler() *TicketRequestHandler {
	return &TicketRequestHandler{
		s: ticket.GetTicketRequestService(),
	}
}

func (h *TicketRequestHandler) FindRequestById(c *fiber.Ctx) error {
	id := c.Params("id")
	response, err := h.s.FindById(id)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.JsonResponse(c, response)
}

func (h *TicketRequestHandler) FindAllRequest(c *fiber.Ctx) error {
	response, err := h.s.FindAll()
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.JsonResponse(c, response)
}

func (h *TicketRequestHandler) FindRequestByMail(c *fiber.Ctx) error {
	mail := c.Params("mail")
	response, err := h.s.FindByMail(mail)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.JsonResponse(c, response)
}

func (h *TicketRequestHandler) FindRequestByStatus(c *fiber.Ctx) error {
	status := c.Params("status")
	response, err := h.s.FindByStatus(status)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.JsonResponse(c, response)
}

func (h *TicketRequestHandler) FindRequestByMailAndStatus(c *fiber.Ctx) error {
	mail := c.Params("mail")
	status := c.Params("status")
	response, err := h.s.FindByMailAndStatus(mail, status)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.JsonResponse(c, response)
}
