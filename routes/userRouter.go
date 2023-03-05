package routes

import (
	"github.com/akgupta-47/auth-gofib/controller"
	"github.com/akgupta-47/auth-gofib/middleware"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App) {
	app.Use(middleware.Authenticate)
	app.Get("/users", controller.GetUsers)
	app.Get("/users/:user_id", controller.GetUser)
}
