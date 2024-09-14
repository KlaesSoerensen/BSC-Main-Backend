package api

import (
	"crypto/rand"
	"encoding/base32"
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"
	"otte_main_backend/src/util"
	"otte_main_backend/src/vitec"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SessionInitiationResponseDTO struct {
	Token string `json:"token"`
}

func applySessionApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Session API] Applying session API")

	//No Auth required
	app.Post("/api/v1/session", func(c *fiber.Ctx) error { return initiateSessionHandler(c, appContext) })

	return nil
}

func initiateSessionHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	var body vitec.SessionInitiationDTO
	//Extract request body
	if err := c.BodyParser(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	var player PlayerDTO
	var idOfPreviousSession int = -1
	//Check if player exists in PlayerDB - if so, all is well
	if err := appContext.PlayerDB.Where("referenceID = ?", body.UserIdentifier).First(&player).Error; err != nil {

		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusInternalServerError)
			middleware.LogRequests(c)
			return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to lookup the user")
		}
		//If the player doesn't exist, check with Vitec
		if crossVerificationError := appContext.VitecIntegration.VerifyUser(&body); crossVerificationError != nil {
			c.Status(fiber.StatusUnauthorized)
			middleware.LogRequests(c)
			return fiber.NewError(fiber.StatusUnauthorized, "Unable to cross-verify")
		}
		//If the user is cross verified but doesn't exist in our system, a new player is created
		player = PlayerDTO{
			ReferenceID: body.UserIdentifier,
			IGN:         body.IGN,
			Sprite:      1,
		}

		if createPlayerError := appContext.PlayerDB.Create(&player).Error; createPlayerError != nil {
			c.Status(fiber.StatusInternalServerError)
			middleware.LogRequests(c)
			return fiber.NewError(fiber.StatusInternalServerError, "Unable to create player")
		}
	} else {
		//If the user exists, grab any earlier session's id
		var session auth.Session
		if err := appContext.PlayerDB.Where("player = ?", player.ID).First(&session).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusInternalServerError)
				middleware.LogRequests(c)
				return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to lookup the session")
			}
		} else {
			//If a session exists, we can just return the token
			idOfPreviousSession = int(session.ID)
		}
	}
	// If this point is reached, the player now exists in the system and is cross-verified (or has been cross-verified before)
	var token, generationErr = generateBase32String(128)
	if generationErr != nil {
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to generate session token")
	}

	var session = auth.Session{
		ID:     uint32(util.Ternary(idOfPreviousSession == -1, 0, idOfPreviousSession)),
		Token:  token,
		Player: player.ID,
	}
	if createSessionError := appContext.PlayerDB.Create(&session).Error; createSessionError != nil {
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to create session")
	}

	c.Status(fiber.StatusOK)
	middleware.LogRequests(c)
	return c.JSON(SessionInitiationResponseDTO{Token: token})
}

func generateBase32String(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(bytes), nil
}
