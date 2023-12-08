package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/controllers"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
)

func UserRoutes(app *fiber.App) {
	app.Post("/user", controllers.CreateUser)
	app.Post("/user/login", controllers.Login)
	app.Get("/users", middlewares.VerifyToken, controllers.GetUsers)
	app.Get("/user/:id", middlewares.VerifyToken, controllers.GetUser)
	app.Put("/user/:id", middlewares.VerifyToken, controllers.UpdateUser)
	app.Delete("/user/:id", middlewares.VerifyToken, controllers.DeleteUser)
}
