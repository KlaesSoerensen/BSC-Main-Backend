package main

import (
	"log"
	api "otte_main_backend/src/api"
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

	var context = meta.CreateApplicationContext(colonyDB, languageDB, playerDB, vitecIntegration, config.GetOr("DEFAULT_DEBUG_HEADER", "DEFAULT-DEBUG-HEADER"))
	app := fiber.New()

	app.Use(cors.New())
	if apiErr := api.ApplyEndpoints(app, context); apiErr != nil {
		panic(apiErr)
	}

	log.Println("[server] Starting server...")
	log.Fatal(doTheTLSThing(servicePort, app))
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

	languageDB, err := db.ConnectLanguageDB()
	if err != nil {
		return nil, nil, nil, err
	}

	playerDB, err := db.ConnectPlayerDB()
	if err != nil {
		return nil, nil, nil, err
	}

	return colonyAssetDB, languageDB, playerDB, nil
}
