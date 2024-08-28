package middleware

import (
	"fmt"
	"otte_main_backend/src/config"

	"github.com/gofiber/fiber/v2"
)

func ApplyAuth(app *fiber.App) error {
	authTokenName, err := config.LoudGet("AUTH_TOKEN_NAME")
	if err != nil {
		return err
	}

	app.Use("/api", func(context *fiber.Ctx) error {
		return naiveCheckForHeaderAuth(context, authTokenName)
	})

	return nil
}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string) error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Status(401).SendString("Unauthorized")
		return fmt.Errorf("Unauthorized")
	}

	return nil
}
