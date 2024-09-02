package api

import (
	"log"
	"otte_main_backend/src/meta"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type LODResponse struct {
	ID          uint32 `json:"id"`
	DetailLevel uint32 `json:"detailLevel"`
	Blob        []byte `json:"blob"`
}

func applyLodApi(app *fiber.App, appContext meta.ApplicationContext) error {
	log.Println("[LOD API] Applying LOD API")

	app.Get("/api/v1/lod/:lodId", func(c *fiber.Ctx) error {
		var idStr = c.Params("lodId")
		id, parsingError := strconv.Atoi(idStr)
		if parsingError != nil {

			c.Status(fiber.StatusBadRequest)
			c.Response().Header.Set(appContext.DDH, "Error occured during parsing of id. Recieved: "+idStr)
			return c.Next()
		}
		var data LODResponse
		if queryErr := appContext.ColonyAssetDB.Table(`LOD`).
			Where(`"LOD".id = ?`, id).
			Select(`"LOD".id, "LOD"."detailLevel", "LOD".blob`).
			First(&data).Error; queryErr != nil {

			c.Status(fiber.StatusInternalServerError)
			c.Response().Header.Set(appContext.DDH, "Error occured during query of LOD. Recieved: "+idStr)
			return c.Next()
		}

		c.Status(fiber.StatusOK)
		return c.JSON(data)
	})

	return nil
}
