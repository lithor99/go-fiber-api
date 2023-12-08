package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/lithor99/go-api-fiber-mysql/controllers"
)

func UploadRoutes(app *fiber.App) {
	app.Post("/upload/single", controllers.UploadSingleFile)
	app.Post("/upload/multi", controllers.UploadMultiFile)
	app.Post("/upload/excel", controllers.UploadExcelFile)
}
