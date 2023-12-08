package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/controllers"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
)

func OrderRoutes(app *fiber.App) {
	app.Post("/order", middlewares.VerifyToken, controllers.CreateOrder)
	app.Get("/orders", middlewares.VerifyToken, controllers.GetOrders)
	app.Get("/order/:id", middlewares.VerifyToken, controllers.GetOrder)
	app.Put("/order/:id", middlewares.VerifyToken, controllers.UpdateOrder)
	app.Delete("/order/:id", middlewares.VerifyToken, controllers.DeleteOrder)
}
