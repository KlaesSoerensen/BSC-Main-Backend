package auth

import (
	"errors"
	"fmt"
	"log"
	"otte_main_backend/src/config"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"
	"otte_main_backend/src/util"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

var errorUnauthorized error = fmt.Errorf("unauthorized")

type CacheEntry[T any] struct {
	Entry     *T
	CreatedAt time.Time
}
type AuthService struct {
	//Valid for a minute before requiring looking up again
	SessionCache util.ConcurrentTypedMap[SessionToken, CacheEntry[Session]]
	//Function to call when authenticating a request
	Method func(c *fiber.Ctx) *fiber.Error
}

var authSingleton *AuthService = &AuthService{
	SessionCache: util.ConcurrentTypedMap[SessionToken, CacheEntry[Session]]{},
	Method:       nil, //Set in InitializeAuth
}

type AuthLevel string

const (
	AuthLevelStrict AuthLevel = "strict"
	AuthLevelNaive  AuthLevel = "naive"
)

func InitializeAuth(appContext *meta.ApplicationContext) (*AuthService, error) {
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
		return nil, fmt.Errorf("Invalid auth level: %s", level)
	}

	return authSingleton, nil
}

func GetSessionByReferenceID(referenceID string, appContext *meta.ApplicationContext) (*Session, error) {
	//First check cache
	authSingleton.SessionCache.Load(SessionToken(referenceID))
	//If not in cache, check DB

	//If in DB, add to cache

	//If not in DB either, return error
	return nil, errorUnauthorized
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

func fullSessionCheckAuth(authService *AuthService, c *fiber.Ctx, appContext *meta.ApplicationContext) *fiber.Error {
	authHeaderContent := string(c.Request().Header.Peek(appContext.AuthTokenName))
	if len(authHeaderContent) == 0 {
		c.Response().Header.Set(appContext.DDH, "Missing auth header, expected "+appContext.AuthTokenName+" to be present")
		return ErrorUnauthorized
	}
	if cacheEntry, exists := authService.SessionCache.Load(SessionToken(authHeaderContent)); exists {
		//If cache entry
		if time.Since(cacheEntry.CreatedAt) < time.Minute {
			//If cache entry is valid (within 1 minute)
			return nil
		}
	}
	//If no cache entry OR cache entry is expired
	var session Session
	var dbErr error
	if dbErr = appContext.PlayerDB.Where("token = ?", authHeaderContent).First(&session).Error; dbErr != nil {
		if !errors.Is(dbErr, gorm.ErrRecordNotFound) {
			log.Println("[AUTH] INTERNAL ERROR: " + dbErr.Error())
		}
		c.Response().Header.Set(appContext.DDH, "Invalid session token")
		return ErrorUnauthorized
	}
	if !IsSessionStillValid(&session) {
		//If the session is expired
		c.Response().Header.Set(appContext.DDH, "Session expired")
		return ErrorUnauthorized
	}
	//Session exists and is valid:
	authService.SessionCache.Store(session.Token, CacheEntry[Session]{Entry: &session, CreatedAt: time.Now()})

	return nil
}

func naiveCheckForHeaderAuth(context *fiber.Ctx, tokenName string, defaultDebugHeader string) *fiber.Error {
	authHeaderContent := context.Request().Header.Peek(tokenName)
	if len(authHeaderContent) == 0 {
		context.Response().Header.Set(defaultDebugHeader, "Missing auth header, expected "+tokenName+" to be present")
		return ErrorUnauthorized
	}

	return nil
}
