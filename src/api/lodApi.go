package api

import (
	"log"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

type LODResponse struct {
	ID          uint32 `json:"id"`
	DetailLevel uint32 `json:"detailLevel"`
	Blob        []byte `json:"blob"`
}

func applyLodApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[LOD API] Applying LOD API")

	return nil
}
