package api

import (
	"errors"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/util"
	"strconv"
	"strings"

	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO's
type AssetResponse struct {
	ID          uint32 `json:"id"`
	UseCase     string `json:"useCase"`
	Type        string `json:"type"`
	Width       uint32 `json:"width"`
	Height      uint32 `json:"height"`
	HasLODs     bool   `json:"hasLODs"`
	Blob        []byte `json:"blob"`
	Alias       string `json:"alias"`
	LODID       uint32 `json:"lod_id"`       // Matches "LOD".id
	DetailLevel string `json:"detail_level"` // Matches "LOD".detailLevel
	LODBlob     []byte `json:"lod_blob"`     // Matches "LOD".blob
}

type MultiAssetResponse []AssetResponse

// Apply the asset API routes
func applyAssetApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Asset API] Applying asset API")

	app.Get("/api/v1/asset/:assetId", auth.PrefixOn(appContext, getAssetByIdHandler))

	app.Get("/api/v1/assets", auth.PrefixOn(appContext, getMultipleAssetsByIds))

	return nil
}

func getMultipleAssetsByIds(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	idsParam := c.Query("ids")
	idStrings := strings.Split(idsParam, ",")

	ids, parseErr := util.ArrayMapTError(idStrings, strconv.Atoi)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Error in parsing asset id "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Error in parsing asset id "+parseErr.Error())
	}

	var assets []AssetResponse
	err := appContext.ColonyAssetDB.
		Table("GraphicalAsset").
		Select(`"GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob`).
		Where(`"GraphicalAsset".id IN ?`, ids).
		Joins(`LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id`).
		Find(&assets).Error

	// This check goes first since the query is for any amount of rows (including 0 as a valid amount)
	if errors.Is(err, gorm.ErrRecordNotFound) || len(assets) == 0 {
		c.Response().Header.Set(appContext.DDH, "No such assets")
		return fiber.NewError(fiber.StatusNotFound, "No such assets")
	}

	if err != nil {
		// Gorm exposes secrets in err when DB is down, so it can't be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(assets)
}

func getAssetByIdHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	var idstr = c.Params("assetId")
	id, parsingError := strconv.Atoi(idstr)

	if parsingError != nil {
		c.Response().Header.Set(appContext.DDH, "Error parsing id "+parsingError.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing id "+parsingError.Error())
	}

	var dto AssetResponse
	err := appContext.ColonyAssetDB.
		Table("GraphicalAsset").
		Select(`"GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob`).
		Where(`"GraphicalAsset".id = ?`, id).
		Joins(`LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id`).
		First(&dto).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No such asset")
			return fiber.NewError(fiber.StatusNotFound, "No such asset")
		}
		// Gorm exposes secrets in err when DB is down, so it can't be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(dto)
}
