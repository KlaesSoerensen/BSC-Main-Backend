package auth

import (
	"fmt"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/openapi"

	"github.com/gofiber/fiber/v2"
)

func ApplyAuth(apiDef *openapi.ApiDefinition) error {

	apiDef.Middleware.AddRisingEdge(func(context *fiber.Ctx, appContext *meta.ApplicationContext) error {
		return naiveCheckForHeaderAuth(context, apiDef.AuthTokenName, apiDef.DDH)
	})

	return nil
}

var tokenName string

func GetNaiveAuth(apiDef *openapi.ApiDefinition) {

}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		context.Status(401).SendString("Unauthorized")
		return fmt.Errorf("Unauthorized")
	}

	return nil
}

func NoAuth(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	return nil
}
