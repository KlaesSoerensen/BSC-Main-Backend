package api

import (
	"otte_main_backend/src/meta"
	"otte_main_backend/src/openapi"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type ServiceStatus struct {
	ColonyDBConnection   bool   `json:"colonyDBStatus"`
	LanguageDBConnection bool   `json:"languageDBStatus"`
	PlayerDBConnection   bool   `json:"playerDBStatus"`
	StatusMessage        string `json:"statusMessage"`
	Timestamp            string `json:"timestamp"`
}

func applyHealthApi(apiDef *openapi.ApiDefinition) {

	apiDef.Add("/api/v1",
		func(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
			return c.SendString("Hello, World!")
		},
		&openapi.EndpointOptions{
			AuthRequired: openapi.AuthNone,
			ResultSet: []openapi.EndpointResult{
				{
					Code: 200,
					Body: reflect.String,
				},
			},
		},
	)
	/*
		app.Get("/api/v1", )

		app.Get("/api/v1/health", func(c *fiber.Ctx) error {
			//Check db connections here
			colonyDBErr := appContext.ColonyAssetDB.Connection(func(tx *gorm.DB) error { return nil })
			languageDBErr := appContext.LanguageDB.Connection(func(tx *gorm.DB) error { return nil })
			playerDBErr := appContext.PlayerDB.Connection(func(tx *gorm.DB) error { return nil })
			var statusMessage string
			if colonyDBErr != nil || languageDBErr != nil || playerDBErr != nil {
				c.Status(fiber.StatusInternalServerError)
				statusMessage = "Error"
			} else {
				c.Status(fiber.StatusOK)
				statusMessage = "OK"
			}
			var status = ServiceStatus{
				StatusMessage:        statusMessage,
				ColonyDBConnection:   colonyDBErr == nil,
				LanguageDBConnection: languageDBErr == nil,
				PlayerDBConnection:   playerDBErr == nil,
				Timestamp:            time.Now().Format(time.RFC3339),
			}
			return c.JSON(status)
		})

		return nil
	*/
}
