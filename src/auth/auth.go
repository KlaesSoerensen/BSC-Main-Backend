package auth

import (
	"fmt"
	"otte_main_backend/src/config"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

func GetNaiveAuth(serviceConstants *config.ServiceConstants) func(*fiber.Ctx, *meta.ApplicationContext) error {
	return func(context *fiber.Ctx, appContext *meta.ApplicationContext) error {
		return naiveCheckForHeaderAuth(context, serviceConstants.AuthToken, serviceConstants.DDH)
	}
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
