package api

import (
	"errors"
	"otte_main_backend/src/meta"
	"strconv"

	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// DTO's
// URL: <base-URL>/asset/<assetId>
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

func applyAssetApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[Asset API] Applying asset API")

	// app.Get(Path, Handler)
	app.Get("/api/v1/asset/:assetId", func(c *fiber.Ctx) error { // Defined handler

		// Fetch id from request params (fiber.Ctx is request context)
		var idstr = c.Params("assetId")

		// Parse to int (look at return type - tuple)
		id, parsingError := strconv.Atoi(idstr)

		// Check if error occured and set status
		if parsingError != nil {
			c.Status(fiber.StatusBadRequest)

			// Move forward in handling process
			c.Response().Header.Set(appContext.DDH, "Error occured during parsing of id. Recieved: "+idstr)
			return c.Next()
		}

		// Variabel named 'dto' of type AssetRresponse (struct)
		var dto AssetResponse
		err := appContext.ColonyAssetDB. // Access ColonyAsset database.
							Table(`GraphicalAsset`).                                                                                       // Define table
							Select(`"GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob`). // Select parameters
							Where(`"GraphicalAsset".id = ?`, id).                                                                          // Select metric
							Joins(`LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id`).                                      // Join foreign key of LOD's
							First(&dto).Error                                                                                              // Scan into dto, error is retuned

		if err != nil { // Check if error occured
			if errors.Is(err, gorm.ErrRecordNotFound) { // If user error
				c.Status(fiber.StatusNotFound)

				return c.Next()
			}

			// If server error
			c.Status(fiber.StatusInternalServerError)
			c.Response().Header.Set(appContext.DDH, err.Error())

			return c.Next()
		}

		c.Status(fiber.StatusOK)
		return c.JSON(dto) // Serialize and return
	})

	return nil
}
