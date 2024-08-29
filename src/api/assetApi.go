package api

import (
	"otte_main_backend/src/meta"

	"log"

	"github.com/gofiber/fiber/v2"
)

type AssetResponse struct {
	ID      uint32          `json:"id"`
	UseCase string          `json:"useCase"`
	Type    string          `json:"type"`
	Width   uint32          `json:"width"`
	Height  uint32          `json:"height"`
	HasLODs bool            `json:"hasLODs"`
	Blob    []byte          `json:"blob"`
	Alias   string          `json:"alias"`
	LODs    []LODDetailsDTO `json:"LODs"`
}

type MultiAssetResponse []AssetResponse

func applyAssetApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Asset API] Applying asset API")

	return nil
}
