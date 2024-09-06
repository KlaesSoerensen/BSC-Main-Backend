package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PlayerInfoResponse represents the data returned for a player's basic information.
type PlayerInfoResponse struct {
	ID                   uint32          `json:"id"`
	IGN                  string          `json:"IGN"`
	Sprite               uint32          `json:"sprite"`
	Achievements         util.PGIntArray `json:"achievements"`
	HasCompletedTutorial bool            `json:"hasCompletedTutorial"`
}

// PlayerPreference represents a single preference item.
type PlayerPreference struct {
	ID              uint32             `json:"id"`
	PreferenceKey   string             `json:"key"`
	ChosenValue     string             `json:"chosenValue"`
	AvailableValues util.PGStringArray `json:"availableValues"` // Use the custom array type
}

// PlayerPreferencesResponse represents the data returned for a player's preferences.
type PlayerPreferencesResponse struct {
	Preferences []PlayerPreference `json:"preferences"`
}

// Apply the Player API routes
func applyPlayerApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Player API] Applying Player API")

	// Route for fetching a single player's info by their ID
	app.Get("/api/v1/player/:playerId", auth.PrefixOn(appContext, getPlayerInfoHandler))

	// Route for fetching a player's preferences by their ID
	app.Get("/api/v1/player/:playerId/preferences", auth.PrefixOn(appContext, getPlayerPreferencesHandler))

	return nil
}

func getPlayerPreferencesHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	// Fetch player preferences and join them with available values
	var preferences []PlayerPreference
	if err := appContext.PlayerDB.
		Table(`PlayerPreference`).
		Where(`"PlayerPreference".player = ?`, playerId).
		Select(`
			"PlayerPreference".id,
			"PlayerPreference"."preferenceKey",
			"PlayerPreference"."chosenValue",
			"AvailablePreference"."availableValues"`).
		Joins(`
			JOIN "AvailablePreference" ON "PlayerPreference"."preferenceKey" = "AvailablePreference"."preferenceKey"`).
		Find(&preferences).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Preferences not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Preferences not found "+err.Error())
		}
		//Gorm be exposing secrets in err when DB is down, so it cant be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Return the preferences in a structured format
	response := PlayerPreferencesResponse{Preferences: preferences}

	c.Status(fiber.StatusOK)
	return c.JSON(response)
}

func getPlayerInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	// Fetch player information from the database
	var player PlayerInfoResponse
	if err := appContext.PlayerDB.
		Table("Player").
		Select(`id, "IGN", sprite, achievements`).
		Where("id = ?", playerId).
		First(&player).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Player not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Player not found "+err.Error())
		}

		//Gorm be exposing secrets in err when DB is down, so it cant be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Compute the 'HasCompletedTutorial' field
	// Achievement id 1 is always the tutorial achievement
	player.HasCompletedTutorial = util.ArrayContains(player.Achievements, 1)

	c.Status(fiber.StatusOK)
	return c.JSON(player)
}
