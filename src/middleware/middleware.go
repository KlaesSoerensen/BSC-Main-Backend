package middleware

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"otte_main_backend/src/meta"
)

func ApplyTo(app *fiber.App, appContext meta.ApplicationContext) error {
	if err := ApplyAuth(app, appContext); err != nil {
		return err
	}

	app.Use(cors.New()) //Default CORS middleware

	app.Use(logRequests)

	return nil
}

func logRequests(c *fiber.Ctx) error {
	log.Println("Request recieved: ", c.Method(), c.Path(), "\t\t at ", time.Now().Format(time.RFC3339), " from ", c.IP(), " \tresponse: ", c.Response().StatusCode())

	return c.Next()
}
