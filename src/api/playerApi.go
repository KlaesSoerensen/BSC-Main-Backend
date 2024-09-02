package api

import (
	"errors"
	"log"
	"otte_main_backend/src/meta"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO's

// PlayerInfoResponse represents the data returned for a player's basic information.
type PlayerInfoResponse struct {
	ID                   uint32   `json:"id"`
	IGN                  string   `json:"IGN"`
	Sprite               uint32   `json:"sprite"`
	Achievements         []uint32 `json:"achievements"`
	HasCompletedTutorial bool     `json:"hasCompletedTutorial"`
}

// PlayerPreference represents a single preference item.
type PlayerPreference struct {
	ID              uint32   `json:"id"`
	Key             string   `json:"key"`
	ChosenValue     string   `json:"chosenValue"`
	AvailableValues []string `json:"availableValues"`
}

// PlayerPreferencesResponse represents the data returned for a player's preferences.
type PlayerPreferencesResponse struct {
	Preferences []PlayerPreference `json:"preferences"`
}

// Apply the Player API routes
func applyPlayerApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Player API] Applying Player API")

	// Route for fetching a single player's info by their ID
	app.Get("/api/v1/player/:playerId", func(c *fiber.Ctx) error {
		playerIdStr := c.Params("playerId")
		playerId, err := strconv.Atoi(playerIdStr)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid player ID"})
		}

		// Fetch player information from the database
		var player PlayerInfoResponse
		err = appContext.ColonyAssetDB.
			Table("Player").
			Select("id, IGN, sprite").
			Where("id = ?", playerId).
			Scan(&player).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{"error": "Player not found"})
			}
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Internal server error"})
		}

		// Compute the 'HasCompletedTutorial' field
		var tutorialCompletedCount int64
		err = appContext.ColonyAssetDB.
			Table("Achievement").
			Where("player = ? AND title = ?", playerId, "Tutorial Completed"). // Assuming "Tutorial Completed" is the title of the tutorial achievement, replaced by 'id (int)'?
			Count(&tutorialCompletedCount).Error
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to check tutorial completion status"})
		}
		player.HasCompletedTutorial = tutorialCompletedCount > 0

		// Fetch achievements for the player
		err = appContext.ColonyAssetDB.
			Table("Achievement").
			Select("id").
			Where("player = ?", playerId).
			Pluck("id", &player.Achievements).Error

		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to fetch achievements"})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(player)
	})

	// Route for fetching a player's preferences by their ID
	app.Get("/api/v1/player/:playerId/preferences", func(c *fiber.Ctx) error {
		playerIdStr := c.Params("playerId")
		playerId, err := strconv.Atoi(playerIdStr)
		if err != nil {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Invalid player ID"})
		}

		// Fetch player preferences from the database
		var preferences []PlayerPreference
		err = appContext.ColonyAssetDB.
			Table("PlayerPreference pp").
			Select("pp.id, pp.key, pp.chosen_value, ap.available_values").
			Joins("JOIN AvailablePreference ap ON pp.key = ap.key").
			Where("pp.player = ?", playerId).
			Scan(&preferences).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{"error": "Preferences not found"})
			}
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Internal server error"})
		}

		// Return the preferences in a structured format
		response := PlayerPreferencesResponse{Preferences: preferences}

		c.Status(fiber.StatusOK)
		return c.JSON(response)
	})

	return nil
}
