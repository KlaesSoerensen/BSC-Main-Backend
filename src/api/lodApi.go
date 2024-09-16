package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
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

	app.Get("/api/v1/lod/:id", auth.PrefixOn(appContext, getLODByIDHandler))

	return nil
}

type LOD struct {
	ID             uint32 `json:"id"`
	DetailLevel    uint32 `json:"detailLevel" gorm:"column:detailLevel"`
	GraphicalAsset uint32 `json:"graphicalAsset" gorm:"column:graphicalAsset"`
	Blob           []byte `json:"blob" gorm:"column:blob"`

	MIMEType string `json:"type" gorm:"column:type"`
}

func (lod *LOD) TableName() string {
	//Gorm would otherwise overwrite the "LOD" to "lo_d" or smth like that
	return "LOD"
}

const detailLevelHeaderName = "URSA-DETAIL-LEVEL"
const assetIDHeaderName = "URSA-ASSET-ID"

func getLODByIDHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	idStr := c.Params("id")
	id, idErr := strconv.ParseUint(idStr, 10, 32)
	if idErr != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid ID: "+idErr.Error())
	}
	var lod LOD
	if dbErr := appContext.ColonyAssetDB.Where("id = ?", id).First(&lod).Error; dbErr != nil {
		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			c.Status(fiber.StatusNotFound)
			c.Response().Header.Set(appContext.DDH, "No such LOD")
			return fiber.NewError(fiber.StatusNotFound, "LOD not found")
		}

		c.Status(fiber.StatusInternalServerError)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	c.Status(fiber.StatusOK)
	c.Response().Header.Set(detailLevelHeaderName, strconv.FormatUint(uint64(lod.DetailLevel), 10))
	c.Response().Header.Set(assetIDHeaderName, strconv.FormatUint(uint64(lod.GraphicalAsset), 10))
	c.Response().Header.SetContentType(lod.MIMEType)
	c.Response().Header.SetContentLength(len(lod.Blob))
	return c.Send(lod.Blob)
}
