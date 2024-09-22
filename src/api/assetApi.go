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
	ID      uint32 `json:"id"`
	UseCase string `json:"useCase" gorm:"column:useCase"`
	Type    string `json:"type"`
	Width   uint32 `json:"width"`
	Height  uint32 `json:"height"`
	/** DEPRECATED see explanation in lodAPI.go
	HasLODs bool   `json:"hasLODs" gorm:"column:hasLODs"`
	Blob    []byte       `json:"blob"`
	*/
	Alias string       `json:"alias"`
	LODs  []LODDetails `json:"LODs" gorm:"foreignKey:GraphicalAsset;references:ID"`
}

type MultiAssetResponse []AssetResponse

// Apply the asset API routes
func applyAssetApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Asset API] Applying asset API")

	app.Get("/api/v1/asset/:assetID/lod/:lodID", auth.PrefixOn(appContext, getLODByAssetAndDetailLevel))

	app.Get("/api/v1/asset/:assetId", auth.PrefixOn(appContext, getAssetByIdHandler))

	app.Get("/api/v1/assets", auth.PrefixOn(appContext, getMultipleAssetsByIds))

	return nil
}

func getLODByAssetAndDetailLevel(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	assetIdStr := c.Params("assetID")
	lodIdStr := c.Params("lodID")

	assetId, assetIdErr := strconv.Atoi(assetIdStr)
	lodId, lodIdErr := strconv.Atoi(lodIdStr)

	if assetIdErr != nil || lodIdErr != nil {
		c.Response().Header.Set(appContext.DDH, "Error parsing asset or LOD id")
		return fiber.NewError(fiber.StatusBadRequest, "Error parsing asset or LOD id")
	}

	var lod LOD
	err := appContext.ColonyAssetDB.
		Where(`"LOD"."graphicalAsset" = ? AND "LOD"."detailLevel" = ?`, assetId, lodId).
		First(&lod).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No such LOD")
			return fiber.NewError(fiber.StatusNotFound, "No such LOD")
		}
		// Gorm exposes secrets in err when DB is down, so it can't be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	c.Status(fiber.StatusOK)
	SetHeadersForLODBlob(c, &lod)
	return c.Send(lod.Blob)
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
		Preload("LODs"). // Preload the LODs field for each asset, can be optimized as this fetches the blob too, which superfluous
		Where(`"GraphicalAsset".id IN ?`, ids).
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

	if len(assets) != len(ids) {
		c.Response().Header.Set(appContext.DDH, "Some assets were not found")
		c.Status(fiber.StatusPartialContent)
	} else {
		c.Status(fiber.StatusOK)
	}
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
		Preload("LODs"). // Preload the LODs field using the foreign key
		Where(`"GraphicalAsset".id = ?`, id).
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
