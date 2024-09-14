package auth

import (
	"fmt"
	"otte_main_backend/src/config"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var errorUnauthorized error = fmt.Errorf("unauthorized")

type the_auth struct {
	//Valid for a minute before requiring looking up again
	SessionCache sync.Map
}

var authSingleton *the_auth = &the_auth{
	SessionCache: sync.Map{},
}

type AuthLevel string

const (
	AuthLevelStrict AuthLevel = "strict"
	AuthLevelNaive  AuthLevel = "naive"
)

func InitializeAuth() error {
	level := config.GetOr("AUTH_LEVEL", "strict")

	return nil
}

// Expands the original handler function's inputs (adding in the appContext) and prefixes an authcheck function.
//
// if the auth check fails, it assures the handler isn't run
//
// Also sets debug header and status code on auth error
func PrefixOn(appContext *meta.ApplicationContext, existingHandler func(c *fiber.Ctx, appContext *meta.ApplicationContext) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := naiveCheckForHeaderAuth(c, appContext.AuthTokenName, appContext.DDH); err != nil {
			c.Status(err.Code)
			middleware.LogRequests(c)
			return err
		}
		handlerErr := existingHandler(c, appContext)
		if fiberErr, ok := handlerErr.(*fiber.Error); ok {
			c.Status(fiberErr.Code)
		} else if handlerErr != nil {
			// Handle other errors
			fmt.Printf("Unknown handler error occurred: %v\n", handlerErr)
		}
		middleware.LogRequests(c)
		return handlerErr
	}
}

var ErrorUnauthorized *fiber.Error = fiber.NewError(401, errorUnauthorized.Error())

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) *fiber.Error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		return ErrorUnauthorized
	}

	return nil
}
