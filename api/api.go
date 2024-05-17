package api

import (
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"ticker-tracer/handler/router"
)

type RequestBody struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func InitServer() {
	app := fiber.New(fiber.Config{
		Immutable: true,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Static("/", "./static")
	RegisterRoutes(app)
	err := app.Listen(":" + "8080")
	if err != nil {
		panic(err)
	}

}
func RegisterRoutes(app *fiber.App) {
	//metrics
	app.Get("/metrics", monitor.New(monitor.Config{Title: "MyService Metrics Page"}))

	tcdd := app.Group("/tcdd")
	router.Tcdd(tcdd)

	// not found
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(404).JSON(fiber.Map{
			"code":    404,
			"message": "404: Not Found",
		})
	})
}
