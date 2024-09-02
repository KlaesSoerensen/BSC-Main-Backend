package main

import (
	"log"
	api "otte_main_backend/src/api"
	"otte_main_backend/src/config"
	db "otte_main_backend/src/database"
	"otte_main_backend/src/meta"
	middleware "otte_main_backend/src/middleware"
	"otte_main_backend/src/openapi"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func main() {
	log.Println("[server] Starting server...")
	if envErr := config.DetectAndApplyENV(); envErr != nil {
		panic(envErr)
	}

	serviceConstants, err := config.NewServiceConstants()
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
	}
	app := fiber.New()
	apiDef, err := openapi.New(app)
	if err != nil {
		panic(err)
	}

	if middlewareErr := middleware.ApplyTo(apiDef); middlewareErr != nil {
		panic(middlewareErr)
	}

	if apiErr := api.ApplyEndpoints(apiDef); apiErr != nil {
		panic(apiErr)
	}
	apiDef.BuildApi(&context)

	log.Fatal(doTheTLSThing(serviceConstants.ServicePort, app))
}

func doTheTLSThing(port int64, app *fiber.App) error {
	//Self signed cert generated following:
	//https://gist.github.com/taoyuan/39d9bc24bafc8cc45663683eae36eb1a
	//See "OTTE Dev Cert Details" file for details
	return app.ListenTLS(
		":"+strconv.FormatInt(port, 10),
		"certs/otte_dev_cert.crt",
		"certs/otte_dev_cert.key")
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
