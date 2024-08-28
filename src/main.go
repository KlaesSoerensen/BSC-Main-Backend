package main

import (
	"log"
	api "otte_main_backend/src/api"
	"otte_main_backend/src/config"
	middleware "otte_main_backend/src/middleware"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("[server] Starting server...")
	if envErr := config.DetectAndApplyENV(); envErr != nil {
		panic(envErr)
	}

	servicePort, err := initServerResources()
	if err != nil {
		panic(err)
	}

	app := fiber.New()
	middlewareErr := middleware.ApplyTo(app)
	if middlewareErr != nil {
		panic(middlewareErr)
	}

	apiErr := api.ApplyEndpoints(app)
	if apiErr != nil {
		panic(apiErr)
	}
	log.Fatal(app.Listen(":" + strconv.FormatInt(servicePort, 10)))
}

type ServicePort = int64

func initServerResources() (ServicePort, error) {
	servicePortStr, err := config.LoudGet("SERVICE_PORT")
	if err != nil {
		return -1, err
	}
	servicePortInt, err := strconv.ParseInt(servicePortStr, 10, 32)
	if err != nil {
		return -1, err
	}

	return servicePortInt, nil
}
