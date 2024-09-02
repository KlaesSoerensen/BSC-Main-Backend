package api

import (
	"errors"
	"otte_main_backend/src/meta"
	"strconv"
	"strings"

	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO's
type AssetResponse struct {
	ID      uint32       `json:"id"`
	UseCase string       `json:"useCase"`
	Type    string       `json:"type"`
	Width   uint32       `json:"width"`
	Height  uint32       `json:"height"`
	HasLODs bool         `json:"hasLODs"`
	Blob    []byte       `json:"blob"`
	Alias   string       `json:"alias"`
	LODs    []LODDetails `json:"LODs" gorm:"foreignKey:GraphicalAsset"`
}

type MultiAssetResponse []AssetResponse

// Apply the asset API routes
func applyAssetApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Asset API] Applying asset API")

	app.Get("/api/v1/asset/:assetId", func(c *fiber.Ctx) error {
		var idstr = c.Params("assetId")
		id, parsingError := strconv.Atoi(idstr)

		if parsingError != nil {
			c.Status(fiber.StatusBadRequest)
			return c.SendString("Invalid ID format")
		}

		if id == 0 {
			c.Status(fiber.StatusNotFound)
			return c.JSON(fiber.Map{"error": "Record not found"})
		}

		var dto AssetResponse
		err := appContext.ColonyAssetDB.
			Table("GraphicalAsset").
			Select(`"GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob`).
			Where(`"GraphicalAsset".id = ?`, id).
			Joins(`LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id`).
			First(&dto).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{"error": "Record not found"})
			}

			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Internal server error"})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(dto)
	})

	app.Get("/api/v1/assets", func(c *fiber.Ctx) error {
		idsParam := c.Query("ids")
		idStrings := strings.Split(idsParam, ",")

		ids := make([]int, len(idStrings))
		for i, idStr := range idStrings {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				c.Status(fiber.StatusBadRequest)
				return c.SendString("Invalid ID format")
			}
			ids[i] = id
		}

		var assets []AssetResponse
		err := appContext.ColonyAssetDB.
			Table("GraphicalAsset").
			Select(`"GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob`).
			Where(`"GraphicalAsset".id IN ?`, ids).
			Joins(`LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id`).
			Find(&assets).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Status(fiber.StatusNotFound)
				return c.JSON(fiber.Map{"error": "Records not found"})
			}

			c.Status(fiber.StatusInternalServerError)
			return c.JSON(fiber.Map{"error": "Internal server error"})
		}

		c.Status(fiber.StatusOK)
		return c.JSON(assets)
	})

	return nil
}
