package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type InternationalizationCatalogue = map[string]string
type AvailableLanguageModel struct {
	ID   uint32 `json:"id" gorm:"column:id;primaryKey"`
	Code string `json:"code" gorm:"column:code"`
	Icon uint32 `json:"icon" gorm:"column:icon"`
}

func (a *AvailableLanguageModel) TableName() string {
	return "AvailableLanguages"
}

type AvailableLanguageDTO struct {
	Code string `json:"code" gorm:"column:code"`
	Icon uint32 `json:"icon" gorm:"column:icon"`
}

type AvailableLanguagesResponseDTO struct {
	Languages []AvailableLanguageDTO `json:"languages"`
}

func (a *AvailableLanguageDTO) TableName() string {
	return "AvailableLanguages"
}

type CatalogueEntry struct {
	Key   string
	Value string
}

func applyCatalog(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Catalog API] Applying catalog API")

	app.Get("/api/v1/catalog/languages", auth.PrefixOn(appContext, getAvailableLanguagesHandler))

	app.Get("/api/v1/catalog/:language", auth.PrefixOn(appContext, getCatalogueForLanguageHandler))

	return nil
}

func getAvailableLanguagesHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	var data []AvailableLanguageDTO
	if dbErr := appContext.LanguageDB.Find(&data).Error; dbErr != nil {
		c.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(&AvailableLanguagesResponseDTO{
		Languages: data,
	})
}

func getCatalogueForLanguageHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	language := c.Params("language")
	keysStr := strings.Trim(c.Query("keys"), " ")
	var isUsingKeys bool = true

	keysSplit := strings.Split(keysStr, ",")
	if len(keysSplit) == 0 || keysStr == "" {
		isUsingKeys = false
	}

	if language == "" {
		c.Response().Header.Set(appContext.DDH, "Language path parameter missing")
		return fiber.NewError(fiber.StatusBadRequest, "Language path parameter missing")
	}
	var data []CatalogueEntry
	var dbErr *fiber.Error
	if isUsingKeys {
		if dbErr = getCatalogueForSpecificKeys(language, keysSplit, appContext, c, &data); dbErr != nil {
			return dbErr
		}
	} else {
		if dbErr = getFullCatalogueForLanguage(language, appContext, c, &data); dbErr != nil {
			return dbErr
		}
	}
	//Statuscode set by getFullCatalogueForLanguage or getCatalogueForSpecificKeys
	// Convert the results into a map
	asMap := make(map[string]string)
	for _, result := range data {
		asMap[result.Key] = result.Value
	}

	c.Status(fiber.StatusOK)
	return c.JSON(asMap)
}

// Also sets headers and status code on ctx if present
func getFullCatalogueForLanguage(language string, appContext *meta.ApplicationContext, ctx *fiber.Ctx, dest *[]CatalogueEntry) *fiber.Error {
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		//This alias is needed as GORM matches on the struct field name, which is dynamic in this case
		Select("key, \"" + language + "\" AS value").
		Find(dest).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		//Gorm be exposing secrets in err when DB is down, so it cant be included in the response
		ctx.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}
	ctx.Status(fiber.StatusOK)
	return nil
}

// Also sets headers and status code on ctx if present
func getCatalogueForSpecificKeys(language string, keys []string, appContext *meta.ApplicationContext, ctx *fiber.Ctx, dest *[]CatalogueEntry) *fiber.Error {
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		//This alias is needed as GORM matches on the struct field name, which is dynamic in this case
		Select("key, \""+language+"\" AS value").
		Where("key IN ?", keys).
		Find(dest).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		//Gorm be exposing secrets in err when DB is down, so it cant be included in the response
		ctx.Response().Header.Set(appContext.DDH, "Internal error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
	}
	ctx.Status(fiber.StatusOK)
	return nil
}
