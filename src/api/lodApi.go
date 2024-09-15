package api

import (
	"log"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

/**
	DELETE ME

	Deprecated due to changes in LOD response. The response is not a binary stream with
	the rest set as headers:
		"URSA-DETAIL-LEVEL": uint32
		"URSA-ASSET-ID": uint32

type LODResponse struct {
	ID          uint32 `json:"id"`
	DetailLevel uint32 `json:"detailLevel"`
	Blob        []byte `json:"blob"`
}
*/

func applyLodApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[LOD API] Applying LOD API")

	return nil
}
