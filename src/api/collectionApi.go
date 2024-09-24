package api

import (
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MinimizedAssetWithTransformDTO struct {
	Width     uint32       `json:"width"`
	Height    uint32       `json:"height"`
	LODs      []LODDetails `json:"LODs" gorm:"column:LODs"`
	Alias     string       `json:"alias"`
	Type      string       `json:"type"`
	Transform TransformDTO `json:"transform"`
}

type AssetCollectionResponse struct {
	ID      uint32                           `json:"id"`
	Name    string                           `json:"name"`
	Entries []MinimizedAssetWithTransformDTO `json:"entries"`
}

func (a *AssetCollectionResponse) TableName() string {
	return "AssetCollection"
}

type RawResult struct {
	AssetCollectionID uint32  `gorm:"column:assetCollectionId"`
	CollectionName    string  `gorm:"column:collectionName"`
	CollectionEntryID uint32  `gorm:"column:collectionEntryId"`
	GraphicalAssetID  uint32  `gorm:"column:graphicalAssetId"`
	Width             int     `gorm:"column:width"`
	Height            int     `gorm:"column:height"`
	Alias             string  `gorm:"column:alias"`
	Type              string  `gorm:"column:type"`
	XOffset           float32 `gorm:"column:xOffset"`
	YOffset           float32 `gorm:"column:yOffset"`
	ZIndex            uint32  `gorm:"column:zIndex"`
	XScale            float32 `gorm:"column:xScale"`
	YScale            float32 `gorm:"column:yScale"`
	LODID             uint32  `gorm:"column:lodId"`
	DetailLevel       int     `gorm:"column:detailLevel"`
}

const collectionQuery = `
SELECT 
	ac.id AS "assetCollectionId",
	ac.name AS "collectionName",
	ce.id AS "collectionEntryId",
	ga.id AS "graphicalAssetId",
	ga.width AS "width",
	ga.height AS "height",
	ga.alias AS "alias",
	ga.type AS "type",
	t."xOffset" AS "xOffset",
	t."yOffset" AS "yOffset",
	t."zIndex" AS "zIndex",
	t."xScale" AS "xScale",
	t."yScale" AS "yScale",
	lod.id AS "lodId",
	lod."detailLevel" AS "detailLevel"
FROM 
	"AssetCollection" ac
JOIN 
	"CollectionEntry" ce ON ce."assetCollection" = ac.id
JOIN 
	"Transform" t ON t.id = ce.transform
JOIN 
	"GraphicalAsset" ga ON ga.id = ce."graphicalAsset"
LEFT JOIN 
	"LOD" lod ON lod."graphicalAsset" = ga.id
WHERE 
	ac.id = ?`

func applyCollectionApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Collection API] Applying collection API")

	app.Get("/api/v1/collection/:collectionId", auth.PrefixOn(appContext, getCollectionByIDHandler))

	return nil
}

func getCollectionByIDHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	collectionId, parseErr := c.ParamsInt("collectionId")
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid collection ID: "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid collection ID: "+parseErr.Error())
	}

	var rawResults []RawResult
	if err := appContext.ColonyAssetDB.Raw(collectionQuery, collectionId).Scan(&rawResults).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.Response().Header.Set(appContext.DDH, "No such collection")
			return fiber.NewError(fiber.StatusNotFound, "No such collection")
		}
		log.Printf("[Collection API] Error retrieving collection: %v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "Error retrieving collection")
	}

	if len(rawResults) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Collection not found"})
	}

	response := transformRawResults(rawResults)
	return c.Status(fiber.StatusOK).JSON(response)
}

func transformRawResults(rawResults []RawResult) *AssetCollectionResponse {
	response := &AssetCollectionResponse{
		ID:      rawResults[0].AssetCollectionID,
		Name:    rawResults[0].CollectionName,
		Entries: []MinimizedAssetWithTransformDTO{},
	}

	entriesMap := make(map[uint32]*MinimizedAssetWithTransformDTO)

	for _, raw := range rawResults {
		entry, exists := entriesMap[raw.CollectionEntryID]
		if !exists {
			entry = &MinimizedAssetWithTransformDTO{
				Width:  uint32(raw.Width),
				Height: uint32(raw.Height),
				Alias:  raw.Alias,
				Type:   raw.Type,
				LODs:   []LODDetails{},
				Transform: TransformDTO{
					XOffset: raw.XOffset,
					YOffset: raw.YOffset,
					ZIndex:  raw.ZIndex,
					XScale:  raw.XScale,
					YScale:  raw.YScale,
				},
			}
			entriesMap[raw.CollectionEntryID] = entry
			response.Entries = append(response.Entries, *entry)
		}

		entry.LODs = append(entry.LODs, LODDetails{
			ID:          raw.LODID,
			DetailLevel: uint32(raw.DetailLevel),
		})

	}

	return response
}
