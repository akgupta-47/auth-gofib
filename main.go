package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/akgupta-47/auth-gofib/db"
	routes "github.com/akgupta-47/auth-gofib/routes"
)

func main() {
	if err := db.ConnectDB(); err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	app := fiber.New()
	// Default middleware config
	app.Use(logger.New())

	routes.AuthRoutes(app)
	routes.UserRoutes(app)

	app.Get("/api", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request!")
	})

	// routes.FRouter(app)
	log.Fatal(app.Listen(":" + port))
}
