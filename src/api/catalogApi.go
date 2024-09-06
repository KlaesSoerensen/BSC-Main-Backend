package api

import (
	"log"
	"otte_main_backend/src/meta"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type InternationalizationCatalogue = map[string]string

func applyCatalog(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Catalog API] Applying catalog API")

	app.Get("/api/v1/catalog/:language", func(c *fiber.Ctx) error {
		language := c.Params("language")
		if language == "" {
			c.Status(fiber.StatusBadRequest)
			return c.JSON(fiber.Map{"error": "Language not provided"})
		}
		keysSplit := strings.Split(c.Query("keys"), ",")
		var data map[string]string
		if len(keysSplit) == 0 {
			appContext.LanguageDB.
				Table("Catalogue").
				Select("key, " + language).
				Scan(&data)
		} else {
			appContext.LanguageDB.
				Table("Catalogue").
				Where("key IN ?", keysSplit).
				Select("key, " + language).Scan(&data)
		}
		c.Status(fiber.StatusOK)
		return c.JSON(data)
	})

	return nil
}
