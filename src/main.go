package main

import (
	"log"
	api "otte_main_backend/src/api"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/config"
	db "otte_main_backend/src/database"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/vitec"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	if _, envErr := config.DetectAndApplyENV(); envErr != nil {
		panic(envErr)
	}

	servicePort, err := initServerResources()
	if err != nil {
		panic(err)
	}

	colonyDB, languageDB, playerDB, dbErr := ConnectDatabases()
	if dbErr != nil {
		panic(dbErr)
	}

	vitecIntegration, integrationErr := vitec.CreateNewVitecIntegration()
	if integrationErr != nil {
		panic(integrationErr)
	}

	context, err := meta.CreateApplicationContext(colonyDB, languageDB, playerDB, vitecIntegration)
	if err != nil {
		panic(err)
	}
	authService, authInitErr := auth.InitializeAuth(context)
	if authInitErr != nil {
		panic(authInitErr)
	}

	app := fiber.New()

	app.Use(cors.New())
	if apiErr := api.ApplyEndpoints(app, context, authService); apiErr != nil {
		panic(apiErr)
	}

	log.Println("[server] Starting server...")
	useTLS := config.GetOr("ENABLE_TLS", "true") == "true"
	if useTLS {
		log.Fatal(doTheTLSThing(servicePort, app))
	} else {
		log.Fatal(listenHTTP(servicePort, app))
	}
}

func doTheTLSThing(port ServicePort, app *fiber.App) error {
	//Self signed cert generated following:
	//https://gist.github.com/taoyuan/39d9bc24bafc8cc45663683eae36eb1a
	//See "OTTE Dev Cert Details" file for details
	return app.ListenTLS(
		":"+strconv.FormatInt(port, 10),
		"certs/otte_dev_cert.crt",
		"certs/otte_dev_cert.key")
}

func listenHTTP(port ServicePort, app *fiber.App) error {
	log.Println("[server] TLS DISABLED.")
	return app.Listen(":" + strconv.FormatInt(port, 10))
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

func ConnectDatabases() (db.ColonyAssetDB, db.LanguageDB, db.PlayerDB, error) {
	colonyAssetDB, err := db.ConnectColonyAssetDB()
	if err != nil {
		return nil, nil, nil, err
	}
	log.Println("[database] Successfully connected to Colony Asset DB")

	languageDB, err := db.ConnectLanguageDB()
	if err != nil {
		return nil, nil, nil, err
	}
	log.Println("[database] Successfully connected to Language DB")

	playerDB, err := db.ConnectPlayerDB()
	if err != nil {
		return nil, nil, nil, err
	}
	log.Println("[database] Successfully connected to Player DB")

	return colonyAssetDB, languageDB, playerDB, nil
}
