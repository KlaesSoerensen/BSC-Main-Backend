package auth

import (
	"fmt"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"

	"github.com/gofiber/fiber/v2"
)

var errorUnauthorized error = fmt.Errorf("Unauthorized")

// Expands the original handler function's inputs (adding in the appContext) and prefixes an authcheck function.
//
// if the auth check fails, it assures the handler isn't run
//
// Also sets debug header and status code on auth error
func PrefixOn(appContext *meta.ApplicationContext, existingHandler func(c *fiber.Ctx, appContext *meta.ApplicationContext) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := naiveCheckForHeaderAuth(c, appContext.AuthTokenName, appContext.DDH); err != nil {
			middleware.LogRequests(c)
			return nil
		}
		handlerErr := existingHandler(c, appContext)
		middleware.LogRequests(c)
		return handlerErr
	}
}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		context.Status(401)
		return errorUnauthorized
	}

	return nil
}
