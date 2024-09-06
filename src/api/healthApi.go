package api

import (
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ServiceStatus struct {
	ColonyDBConnection   bool   `json:"colonyDBStatus"`
	LanguageDBConnection bool   `json:"languageDBStatus"`
	PlayerDBConnection   bool   `json:"playerDBStatus"`
	StatusMessage        string `json:"statusMessage"`
	Timestamp            string `json:"timestamp"`
}

func rootHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	return c.Status(fiber.StatusOK).SendString("You've reached the backend.")
}

func healthRouteHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
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
		StatusMessage:        statusMessage,
		ColonyDBConnection:   colonyDBErr == nil,
		LanguageDBConnection: languageDBErr == nil,
		PlayerDBConnection:   playerDBErr == nil,
		Timestamp:            time.Now().Format(time.RFC3339),
	}
	return c.JSON(status)
}

func applyHealthApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Health API] Applying health API")

	app.Get("/api/v1", auth.PrefixOn(appContext, rootHandler))

	app.Get("/api/v1/health", auth.PrefixOn(appContext, healthRouteHandler))

	return nil
}
