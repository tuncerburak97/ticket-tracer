package router

import (
	"github.com/gofiber/fiber/v2"
	"ticker-tracer/api/handler"
)

func TicketRequest(router fiber.Router) {
	var ticketRequestHandler = handler.NewTicketRequestHandler()
	router.Get("/:id", ticketRequestHandler.FindRequestById)
	router.Get("", ticketRequestHandler.FindAllRequest)
	router.Get("/mail/:mail", ticketRequestHandler.FindRequestByMail)
	router.Get("/status/:status", ticketRequestHandler.FindRequestByStatus)
	router.Get("/mail/:mail/status/:status", ticketRequestHandler.FindRequestByMailAndStatus)
}
