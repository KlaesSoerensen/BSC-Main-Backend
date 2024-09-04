package api

import (
	"database/sql/driver"
	"errors"
	"log"
	"otte_main_backend/src/meta"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// StringArray is a custom type to handle PostgreSQL arrays.
type StringArray []string

// Scan implements the Scanner interface for StringArray.
func (a *StringArray) Scan(value interface{}) error {
	return pq.Array(a).Scan(value)
}

// Value implements the Valuer interface for StringArray.
func (a StringArray) Value() (driver.Value, error) {
	return pq.Array(a).Value()
}

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
	ID              uint32      `json:"id"`
	Key             string      `json:"key"`
	ChosenValue     string      `json:"chosenValue"`
	AvailableValues StringArray `json:"availableValues"` // Use the custom array type
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
		err = appContext.PlayerDB.
			Table("Player").
			Select(`id, "IGN", sprite`).
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
		err = appContext.PlayerDB.
			Table("Achievement").
			Where("player = ? AND title = ?", playerId, "Tutorial Completed"). // Assuming "Tutorial Completed" is the title of the tutorial achievement
			Count(&tutorialCompletedCount).Error
		if err != nil {
			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Failed to check tutorial completion status"})
		}
		player.HasCompletedTutorial = tutorialCompletedCount > 0

		// Fetch achievements for the player
		err = appContext.PlayerDB.
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

		// Fetch player preferences and join them with available values
		var preferences []PlayerPreference
		err = appContext.PlayerDB.
			Table(`PlayerPreference`).
			Select(`
            "PlayerPreference".id,
            "PlayerPreference"."preferenceKey",
            "PlayerPreference"."chosenValue",
            "AvailablePreference"."availableValues"
        `).
			Joins(`
            JOIN "AvailablePreference" ON "PlayerPreference"."preferenceKey" = "AvailablePreference"."preferenceKey"
        `).
			Where(`"PlayerPreference".player = ?`, playerId).
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
