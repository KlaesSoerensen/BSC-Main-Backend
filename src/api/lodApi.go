package api

import (
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

func applyLodApi(app *fiber.App, appContext meta.ApplicationContext) error {
	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	return nil
}
