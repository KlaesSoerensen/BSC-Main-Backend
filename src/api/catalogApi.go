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

func applyCatalog(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Catalog API] Applying catalog API")

	app.Get("/api/v1/catalog/:language", auth.PrefixOn(appContext, getCatalogueForLanguageHandler))

	return nil
}

type CatalogueEntry struct {
	Key   string
	Value string
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
	json := make(map[string]string)
	for _, result := range data {
		json[result.Key] = result.Value
	}

	c.Status(fiber.StatusOK)
	return c.JSON(json)
}

// Also sets headers and status code on ctx if present
func getFullCatalogueForLanguage(language string, appContext *meta.ApplicationContext, ctx *fiber.Ctx, dest *[]CatalogueEntry) *fiber.Error {
	log.Printf("[delete me] Executing query for language: %s", language)
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		//This alias is needed as GORM matches on the struct field name, which is dynamic in this case
		Select("key, \"" + language + "\" AS value").
		Scan(dest).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		ctx.Response().Header.Set(appContext.DDH, "Internal error "+dbErr.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error "+dbErr.Error())
	}
	ctx.Status(fiber.StatusOK)
	return nil
}

// Also sets headers and status code on ctx if present
func getCatalogueForSpecificKeys(language string, keys []string, appContext *meta.ApplicationContext, ctx *fiber.Ctx, dest *[]CatalogueEntry) *fiber.Error {
	log.Printf("[delete me] Executing query for language: %s, keys: %v", language, keys)
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		//This alias is needed as GORM matches on the struct field name, which is dynamic in this case
		Select("key, \""+language+"\" AS value").
		Where("key IN ?", keys).
		Scan(dest).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		ctx.Response().Header.Set(appContext.DDH, "Internal error "+dbErr.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal error "+dbErr.Error())
	}
	ctx.Status(fiber.StatusOK)
	return nil
}
