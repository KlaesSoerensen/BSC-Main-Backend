package api

import (
	"log"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

func applyCatalog(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Catalog API] Applying catalog API")
	return nil
}
