package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApplyTo(app *fiber.App) error {
	app.Use(cors.New()) //Default CORS middleware
	app.Use(logRequests)
	return ApplyAuth(app)
}

func logRequests(c *fiber.Ctx) error {
	log.Println("Request recieved: ", c.Method(), c.Path(), "\t\t at ", time.Now().Format(time.RFC3339), " from ", c.IP())
	return c.Next()
}
