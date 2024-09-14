package api

import (
	"log"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"

	"github.com/gofiber/fiber/v2"
)

type SessionInitiationDTO struct {
	UserIdentifier      string `json:"userIdentifier"`
	CurrentSessionToken string `json:"currentSessionToken"`
}

func applySessionApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Session API] Applying session API")

	//No Auth required
	app.Post("/api/v1/session", func(c *fiber.Ctx) error { return initiateSessionHandler(c, appContext) })

	return nil
}

func initiateSessionHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	var body SessionInitiationDTO
	//Extract request body
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	var player PlayerDTO
	//Check if player exists in PlayerDB - if so, all is well
	if err := appContext.PlayerDB.Where("id = ?", body.UserIdentifier).First(&player).Error; err != nil {

	}

	//If not, check with Vitec

	c.Status(fiber.StatusNotImplemented)
	middleware.LogRequests(c)
	return fiber.NewError(fiber.StatusNotImplemented, "Not implemented")
}
