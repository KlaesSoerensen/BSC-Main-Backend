package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ApplyTo(app *fiber.App) error {
	app.Use(cors.New()) //Default CORS middleware
	return ApplyAuth(app)
}
