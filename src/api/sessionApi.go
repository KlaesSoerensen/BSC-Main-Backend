package api

import (
	"log"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"

	"github.com/gofiber/fiber/v2"
)

func applySessionApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Session API] Applying session API")

	//No Auth required
	app.Post("/api/v1/session", func(c *fiber.Ctx) error { return initiateSessionHandler(c, appContext) })

	return nil
}

func initiateSessionHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {

	c.Status(fiber.StatusInternalServerError)
	middleware.LogRequests(c)
	return fiber.NewError(fiber.StatusNotImplemented, "Not implemented")
}
