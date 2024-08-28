package api

import (
	"log"
	"otte_main_backend/src/meta"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ServiceStatus struct {
	ColonyDBConnectionError   bool   `json:"colonyDBStatus"`
	LanguageDBConnectionError bool   `json:"languageDBStatus"`
	PlayerDBConnectionError   bool   `json:"playerDBStatus"`
	StatusMessage             string `json:"statusMessage"`
	Timestamp                 string `json:"timestamp"`
}

func applyHealthApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Health API] Applying health API")

	app.Get("/api/v1", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/api/v1/health", func(c *fiber.Ctx) error {
		//Check db connections here
		colonyDBErr := appContext.ColonyAssetDB.Connection(func(tx *gorm.DB) error { return nil })
		languageDBErr := appContext.LanguageDB.Connection(func(tx *gorm.DB) error { return nil })
		playerDBErr := appContext.PlayerDB.Connection(func(tx *gorm.DB) error { return nil })
		var statusMessage string
		if colonyDBErr != nil || languageDBErr != nil || playerDBErr != nil {
			c.Status(fiber.StatusInternalServerError)
			statusMessage = "Error"
		} else {
			c.Status(fiber.StatusOK)
			statusMessage = "OK"
		}
		var status = ServiceStatus{
			StatusMessage:             statusMessage,
			ColonyDBConnectionError:   colonyDBErr != nil,
			LanguageDBConnectionError: languageDBErr != nil,
			PlayerDBConnectionError:   playerDBErr != nil,
			Timestamp:                 time.Now().Format(time.RFC3339),
		}
		return c.JSON(status)
	})

	return nil
}
