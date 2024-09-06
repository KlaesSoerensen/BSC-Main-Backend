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
	var data map[string]string = make(map[string]string)
	var dbErr *fiber.Error
	if isUsingKeys {
		if data, dbErr = getCatalogueForSpecificKeys(language, keysSplit, appContext, c); dbErr != nil {
			return dbErr
		}
	} else {
		if data, dbErr = getFullCatalogueForLanguage(language, appContext, c); dbErr != nil {
			return dbErr
		}
	}
	//Statuscode set by getFullCatalogueForLanguage or getCatalogueForSpecificKeys
	return c.JSON(data)
}

// Also sets headers and status code on ctx if present
func getFullCatalogueForLanguage(language string, appContext *meta.ApplicationContext, ctx *fiber.Ctx) (InternationalizationCatalogue, *fiber.Error) {
	var data map[string]string
	log.Printf("[delete me] Executing query for language: %s", language)
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		Select("key, \"" + language + "\"").
		Scan(&data).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return nil, fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		ctx.Response().Header.Set(appContext.DDH, "Internal error "+dbErr.Error())
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error "+dbErr.Error())
	}
	ctx.Status(fiber.StatusOK)
	return data, nil
}

// Also sets headers and status code on ctx if present
func getCatalogueForSpecificKeys(language string, keys []string, appContext *meta.ApplicationContext, ctx *fiber.Ctx) (InternationalizationCatalogue, *fiber.Error) {
	var data map[string]string
	log.Printf("[deleteme] Executing query for language: %s, keys: %v", language, keys)
	if dbErr := appContext.LanguageDB.
		Table("Catalogue").
		Select("key, \""+language+"\"").
		Where("key IN ?", keys).
		Scan(&data).Error; dbErr != nil {

		if errors.Is(dbErr, gorm.ErrRecordNotFound) {
			ctx.Response().Header.Set(appContext.DDH, "No such language catalogue found")
			return nil, fiber.NewError(fiber.StatusNotFound, "No such language catalogue found")
		}

		ctx.Response().Header.Set(appContext.DDH, "Internal error "+dbErr.Error())
		return nil, fiber.NewError(fiber.StatusInternalServerError, "Internal error "+dbErr.Error())
	}
	ctx.Status(fiber.StatusOK)
	return data, nil
}
