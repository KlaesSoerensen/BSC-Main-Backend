package auth

import (
	"fmt"
	"log"
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
	Method       func(c *fiber.Ctx) *fiber.Error
}

var authSingleton *the_auth = &the_auth{
	SessionCache: sync.Map{},
}

type AuthLevel string

const (
	AuthLevelStrict AuthLevel = "strict"
	AuthLevelNaive  AuthLevel = "naive"
)

func InitializeAuth(appContext *meta.ApplicationContext) error {
	level := config.GetOr("INTERNAL_AUTH_LEVEL", "strict")

	switch AuthLevel(level) {
	case AuthLevelStrict:
		authSingleton.Method = func(c *fiber.Ctx) *fiber.Error {
			return fullSessionCheckAuth(authSingleton, c, appContext)
		}
		log.Println("[AUTH] Level set to strict")
	case AuthLevelNaive:
		authSingleton.Method = func(c *fiber.Ctx) *fiber.Error {
			return naiveCheckForHeaderAuth(c, appContext.AuthTokenName, appContext.DDH)
		}
		log.Println("[AUTH] Level set to naive")
	default:
		return fmt.Errorf("Invalid auth level: %s", level)
	}

	return nil
}

// Expands the original handler function's inputs (adding in the appContext) and prefixes an authcheck function.
//
// if the auth check fails, it assures the handler isn't run
//
// Also sets debug header and status code on auth error
func PrefixOn(appContext *meta.ApplicationContext, existingHandler func(c *fiber.Ctx, appContext *meta.ApplicationContext) error) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if err := authSingleton.Method(c); err != nil {
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

func fullSessionCheckAuth(authService *the_auth, c *fiber.Ctx, appContext *meta.ApplicationContext) *fiber.Error {
	//authHeaderContent := string(c.Request().Header.Peek(appContext.AuthTokenName))

	return fiber.NewError(501, "Not implemented")
}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) *fiber.Error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		return ErrorUnauthorized
	}

	return nil
}
