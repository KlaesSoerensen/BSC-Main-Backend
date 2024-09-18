package api

import (
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

func applySessionApi(app *fiber.App, appContext *meta.ApplicationContext, authService *auth.AuthService) error {
	log.Println("[Session API] Applying session API")

	//No Auth required
	app.Post("/api/v1/session", func(c *fiber.Ctx) error { return initiateSessionHandler(c, appContext, authService) })

	return nil
}

func initiateSessionHandler(c *fiber.Ctx, appContext *meta.ApplicationContext, authService *auth.AuthService) error {
	var body vitec.SessionInitiationDTO
	//Extract request body
	if err := c.BodyParser(&body); err != nil {
		c.Status(fiber.StatusBadRequest)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	//First of all, check if the auth header is present. If so, lookup in the cache and return if found
	existingAuthHeader := c.Request().Header.Peek(appContext.AuthTokenName)
	if existingAuthHeader != nil {
		if cacheEntry, exists := authService.SessionCache.Load(auth.SessionToken(existingAuthHeader)); exists && auth.IsSessionStillValid(cacheEntry.Entry) {

			//update last checkin async
			go auth.UpdateLastPlayerCheckin(cacheEntry.Entry, appContext)
			//No need to check for cache entry expiry here, as the cache is checked for expiry on subsequent request
			//Even if the cache is expired, the session is still valid
			c.Status(fiber.StatusOK)
			middleware.LogRequests(c)
			return c.JSON(SessionInitiationResponseDTO{Token: string(cacheEntry.Entry.Token)})
		}
	}

	var player PlayerDTO
	var idOfPreviousSession int = -1
	//Check if player exists in PlayerDB - if so, all is well
	if err := appContext.PlayerDB.Where(`"referenceID" = ?`, body.UserIdentifier).First(&player).Error; err != nil {
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
			ReferenceID:  body.UserIdentifier,
			IGN:          body.IGN,
			Sprite:       1,
			Achievements: util.PGIntArray{},
		}

		if createPlayerError := appContext.PlayerDB.Create(&player).Error; createPlayerError != nil {
			c.Status(fiber.StatusInternalServerError)
			middleware.LogRequests(c)
			return fiber.NewError(fiber.StatusInternalServerError, "Unable to create player")
		}
	} else {
		//If the user exists, grab any earlier session
		var session auth.Session
		if err := appContext.PlayerDB.Where("player = ?", player.ID).First(&session).Error; err == nil {
			//A session exists, if still valid, just return the token
			idOfPreviousSession = int(session.ID)
			if auth.IsSessionStillValid(&session) {
				//update last checkin async
				go auth.UpdateLastPlayerCheckin(&session, appContext)
				c.Status(fiber.StatusOK)
				middleware.LogRequests(c)
				return c.JSON(SessionInitiationResponseDTO{Token: string(session.Token)})
			}
		} else {
			//If the error is not the expected "record not found", something might have gone horribly wrong, so log it
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusInternalServerError)
				middleware.LogRequests(c)
				return fiber.NewError(fiber.StatusInternalServerError, "Something went wrong when trying to lookup the session")
			}
		}
	}

	// If this point is reached, the player now exists in the system and is cross-verified (or has been cross-verified before)
	session, sessionErr := auth.CreateSessionForPlayer(uint32(player.ID), appContext, idOfPreviousSession, authService)
	if sessionErr != nil {
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Unable to initialize session: "+sessionErr.Error())
	}

	c.Status(fiber.StatusOK)
	middleware.LogRequests(c)
	return c.JSON(SessionInitiationResponseDTO{Token: string(session.Token)})
}
