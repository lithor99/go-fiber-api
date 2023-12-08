package main

import (
	"github.com/lithor99/go-api-fiber-mysql/configs"
	"github.com/lithor99/go-api-fiber-mysql/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()

	//root url
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(&fiber.Map{"data": "Welcome To API By Go Fiber and MySql"})
	})

	//run database
	configs.ConnectDB()
	//use routes
	routes.UserRoutes(app)
	routes.ProductRoutes(app)
	routes.UploadRoutes(app)
	routes.OrderRoutes(app)
	app.Listen(":8000")
}
