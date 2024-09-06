package api

import (
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

type MinimizedAssetWithTransformDTO struct {
	HasLODs   bool         `json:"hasLODs"`
	Width     uint32       `json:"width"`
	Height    uint32       `json:"height"`
	LODs      []LODDetails `json:"LODs"`
	Transform TransformDTO `json:"transform"`
}

type AssetCollectionResponse struct {
	ID      uint32                           `json:"id"`
	Name    string                           `json:"name"`
	Entries []MinimizedAssetWithTransformDTO `json:"entries"`
}

type RawResult struct {
	AssetCollectionID uint32  `json:"assetCollectionId"`
	CollectionName    string  `json:"collectionName"`
	CollectionEntryID uint32  `json:"collectionEntryId"`
	GraphicalAssetID  uint32  `json:"graphicalAssetId"`
	HasLODs           bool    `json:"hasLODs"`
	Width             int     `json:"width"`
	Height            int     `json:"height"`
	XOffset           float32 `json:"xOffset"`
	YOffset           float32 `json:"yOffset"`
	ZIndex            uint32  `json:"zIndex"`
	XScale            float32 `json:"xScale"`
	YScale            float32 `json:"yScale"`
	LODID             *uint32 `json:"lodId"`       // Nullable because it can be null
	DetailLevel       *int    `json:"detailLevel"` // Nullable because it can be null
}

const collectionQuery = `
SELECT 
	ac.id AS "assetCollectionId",
	ac.name AS "collectionName",
	ce.id AS "collectionEntryId",
	ga.id AS "graphicalAssetId",
	ga.hasLOD AS "hasLODs",
	ga.width AS "width",
	ga.height AS "height",
	t.xOffset AS "xOffset",
	t.yOffset AS "yOffset",
	t.zIndex AS "zIndex",
	t.xScale AS "xScale",
	t.yScale AS "yScale",
	lod.id AS "lodId",
	lod.detailLevel AS "detailLevel"
FROM 
	"AssetCollection" ac
JOIN 
	"CollectionEntry" ce ON ce.assetCollection = ac.id
JOIN 
	"Transform" t ON t.id = ce.transform
JOIN 
	"GraphicalAsset" ga ON ga.id = ce.graphicalAsset
LEFT JOIN 
	"LOD" lod ON lod.graphicalAsset = ga.id
WHERE 
	ac.id = ?`

func applyCollectionApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Collection API] Applying collection API")

	app.Get("/api/v1/collection/:collectionId", auth.PrefixOn(appContext, getCollectionByIDHandler))

	return nil
}

func getCollectionByIDHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	collectionId := c.Params("collectionId")
	if collectionId == "" {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"error": "Collection ID not provided"})
	}

	var rawResults []RawResult
	if err := appContext.ColonyAssetDB.Raw(collectionQuery, collectionId).Scan(&rawResults).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		c.Response().Header.Set(appContext.DDH, err.Error())
		return c.Next()
	}

	// Transform the raw results into the structured response
	response := &AssetCollectionResponse{
		ID:      rawResults[0].AssetCollectionID,
		Name:    rawResults[0].CollectionName,
		Entries: []MinimizedAssetWithTransformDTO{},
	}

	entriesMap := make(map[uint32]*MinimizedAssetWithTransformDTO)

	for _, raw := range rawResults {
		// Check if we've already added this CollectionEntry
		if _, exists := entriesMap[raw.CollectionEntryID]; !exists {
			entriesMap[raw.CollectionEntryID] = &MinimizedAssetWithTransformDTO{
				HasLODs: raw.HasLODs,
				Width:   uint32(raw.Width),
				Height:  uint32(raw.Height),
				LODs:    []LODDetails{},
				Transform: TransformDTO{
					XOffset: raw.XOffset,
					YOffset: raw.YOffset,
					ZIndex:  raw.ZIndex,
					XScale:  raw.XScale,
					YScale:  raw.YScale,
				},
			}
			// Add to the response entries
			response.Entries = append(response.Entries, *entriesMap[raw.CollectionEntryID])
		}

		// Add LOD details if present
		if raw.LODID != nil {
			lod := LODDetails{
				ID:          *raw.LODID,
				DetailLevel: uint32(*raw.DetailLevel),
			}
			entriesMap[raw.CollectionEntryID].LODs = append(entriesMap[raw.CollectionEntryID].LODs, lod)
		}
	}

	c.Status(fiber.StatusOK)
	return c.JSON(response)
}
