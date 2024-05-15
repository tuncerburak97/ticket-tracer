package tcdd

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	"ticker-tracer/service/tcdd/model"
	utils "ticker-tracer/util/response"
)
import service "ticker-tracer/service/tcdd"

type TcddHandlerInterface interface {
	AddSearchRequest(c *fiber.Ctx) error
	LoadStations(c *fiber.Ctx) error
	QueryTrain(c *fiber.Ctx) error
}

type TcddHandler struct {
	s *service.TccdService
}

func NewFoodRecipeHandler() *TcddHandler {
	return &TcddHandler{s: service.NewTcddService()}
}

func (h *TcddHandler) AddSearchRequest(c *fiber.Ctx) error {
	var req model.SearchTrainRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.FailResponse(c, err.Error())
	}
	recipe, err := h.s.AddSearchRequest(&req)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.DataResponseCreated(c, recipe)
}

func (h *TcddHandler) LoadStations(c *fiber.Ctx) error {
	stations, err := h.s.LoadStations()
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.DataResponse(c, stations)
}

func (h *TcddHandler) QueryTrain(c *fiber.Ctx) error {
	var req model.QueryTrainRequest
	if err := c.BodyParser(&req); err != nil {
		return utils.FailResponse(c, err.Error())
	}

	recipe, err := h.s.QueryTrain(&req)
	if err != nil {
		return utils.FailResponse(c, err.Error())
	}
	return utils.DataResponse(c, recipe)

}
