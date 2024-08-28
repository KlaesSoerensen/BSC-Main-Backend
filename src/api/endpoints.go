package api

import (
	"github.com/gofiber/fiber/v2"
)

func ApplyEndpoints(app *fiber.App) error {
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return nil
}
