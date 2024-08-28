package main

import (
	"log"
	api "otte_main_backend/src/api"
	"otte_main_backend/src/config"
	db "otte_main_backend/src/database"
	"otte_main_backend/src/meta"
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

	colonyDB, languageDB, playerDB, dbErr := ConnectDatabases()
	if dbErr != nil {
		panic(dbErr)
	}
	var context = meta.ApplicationContext{
		ColonyAssetDB: colonyDB,
		LanguageDB:    languageDB,
		PlayerDB:      playerDB,
		DDH:           config.GetOr("DEFAULT_DEBUG_HEADER", "DEFAULT-DEBUG-HEADER"),
	}

	app := fiber.New()
	if middlewareErr := middleware.ApplyTo(app, context); middlewareErr != nil {
		panic(middlewareErr)
	}

	if apiErr := api.ApplyEndpoints(app, context); apiErr != nil {
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
