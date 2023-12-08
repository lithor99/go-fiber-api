package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/controllers"
	"github.com/lithor99/go-api-fiber-mysql/middlewares"
)

func ProductRoutes(app *fiber.App) {
	app.Post("/product", middlewares.VerifyToken, controllers.CreateProduct)
	app.Post("/product/image/:id", middlewares.VerifyToken, controllers.UploadProductImage)
	app.Get("/products", middlewares.VerifyToken, controllers.GetProducts)
	app.Get("/product/:id", middlewares.VerifyToken, controllers.GetProduct)
	app.Put("/product/:id", middlewares.VerifyToken, controllers.UpdateProduct)
	app.Delete("/product/:id", middlewares.VerifyToken, controllers.DeleteProduct)
}
