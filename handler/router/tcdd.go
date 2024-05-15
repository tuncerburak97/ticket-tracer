package router

import "github.com/gofiber/fiber/v2"
import handler "ticker-tracer/handler/tcdd"

func Tcdd(router fiber.Router) {
	var recipeHandler = handler.NewFoodRecipeHandler()
	router.Post("/add", recipeHandler.AddSearchRequest)
	router.Get("/load", recipeHandler.LoadStations)
	router.Post("/query", recipeHandler.QueryTrain)
}
