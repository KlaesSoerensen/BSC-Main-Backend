package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"otte_main_backend/src/meta"
	"otte_main_backend/src/openapi"
)

func ApplyTo(apiDef *openapi.ApiDefinition) error {
	apiDef.Middleware.Add(cors.New()) //Default CORS middleware

	apiDef.Middleware.AddFallingEdge(func(ctx *fiber.Ctx, appContext *meta.ApplicationContext) error {
		return logRequests(ctx)
	})

	return nil
}

func logRequests(c *fiber.Ctx) error {
	log.Println("Request recieved: ", c.Method(), c.Path(), "\t\t at ", time.Now().Format(time.RFC3339), " from ", c.IP(), " \tresponse: ", c.Response().StatusCode())

	return c.Next()
}
