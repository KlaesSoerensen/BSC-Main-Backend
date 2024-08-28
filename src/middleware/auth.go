package middleware

import (
	"log"
	"otte_main_backend/src/config"

	"github.com/gofiber/fiber/v2"
)

func ApplyAuth(app *fiber.App) error {
	authTokenName, err := config.LoudGet("AUTH_TOKEN_NAME")
	if err != nil {
		return err
	}
	defaultDebugHeader, err := config.LoudGet("DEFAULT_DEBUG_HEADER")
	if err != nil {
		return err
	}
	app.Use(func(context *fiber.Ctx) error {
		return naiveCheckForHeaderAuth(context, authTokenName, defaultDebugHeader)
	})

	return nil
}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	log.Println("Auth header content: ", string(authHeaderContent))
	if len(authHeaderContent) <= 6 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		return context.Status(401).SendString("Unauthorized")
	}

	return context.Next()
}
