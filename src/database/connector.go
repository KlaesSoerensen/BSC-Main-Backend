package database

import (
	"log"
	"otte_main_backend/src/config"
	"strconv"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func ConnectPlayerDB() (PlayerDB, error) {
	port, portErr := strconv.ParseUint(config.Get("PLAYER_DB_PORT"), 10, 32)
	timeout, timeoutErr := strconv.Atoi(config.Get("DB_MAX_TIMEOUT"))

	if portErr != nil {
		log.Println("Error parsing player db port value from environment")
	}
	if timeoutErr != nil {
		log.Println("Error parsing db connection timeout value from environment")
	}

	// Connection URL to connect to Postgres Database
	dsn := DBDSN{
		Host:     config.Get("PLAYER_DB_HOST"),
		Port:     port,
		Username: config.Get("PLAYER_DB_USERNAME"),
		Password: config.Get("PLAYER_DB_PASSWORD"),
		Database: config.Get("PLAYER_DB_NAME"),
		SSLMode:  "disable",
	}

	return attemptConnectionWithinTimeout(timeout, dsn, false)
}

func ConnectLanguageDB() (LanguageDB, error) {
	port, portErr := strconv.ParseUint(config.Get("LANGUAGE_DB_PORT"), 10, 32)
	timeout, timeoutErr := strconv.Atoi(config.Get("DB_MAX_TIMEOUT"))
	if portErr != nil {
		log.Println("Error parsing language db port value from environment")
	}
	if timeoutErr != nil {
		log.Println("Error parsing db connection timeout value from environment")
	}
	dsn := DBDSN{
		Host:     config.Get("LANGUAGE_DB_HOST"),
		Port:     port,
		Username: config.Get("LANGUAGE_DB_USERNAME"),
		Password: config.Get("LANGUAGE_DB_PASSWORD"),
		Database: config.Get("LANGUAGE_DB_NAME"),
		SSLMode:  "disable",
	}

	return attemptConnectionWithinTimeout(timeout, dsn, false)
}

func ConnectColonyAssetDB() (LanguageDB, error) {
	port, portErr := strconv.ParseUint(config.Get("COLONY_ASSET_DB_PORT"), 10, 32)
	timeout, timeoutErr := strconv.Atoi(config.Get("DB_MAX_TIMEOUT"))

	if portErr != nil {
		log.Println("Error parsing colony and asset db port value from environment")
	}
	if timeoutErr != nil {
		log.Println("Error parsing db connection timeout value from environment")
	}
	dsn := DBDSN{
		Host:     config.Get("COLONY_ASSET_DB_HOST"),
		Port:     port,
		Username: config.Get("COLONY_ASSET_DB_USERNAME"),
		Password: config.Get("COLONY_ASSET_DB_PASSWORD"),
		Database: config.Get("COLONY_ASSET_DB_NAME"),
		SSLMode:  "disable",
	}

	return attemptConnectionWithinTimeout(timeout, dsn, false)
}

func attemptConnectionWithinTimeout(timeout int, dsn DBDSN, silencePlease bool) (*gorm.DB, error) {
	log.Printf("[database] Trying to establish connection to %s within: %d seconds. \n", dsn.Database, timeout)
	log.Println("[database] Using dsn: " + dsn.SafeString())

	var err error
	var db *gorm.DB

	var gormConfig = &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}
	if silencePlease {
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	for attemptNum := 1; attemptNum < timeout; attemptNum++ {

		var expectedError error

		log.Printf("[database] %s Attempt %d/%d\n", dsn.Database, attemptNum, timeout)

		// instead of default table name "<>+s" the table name will be "<>"
		// if succesfully connected
		// if last attempt
		db, expectedError = gorm.Open(postgres.Open(dsn.FullString()), gormConfig)

		if expectedError == nil {
			break
		}

		//On last try, return the error
		if attemptNum == timeout-1 {
			err = expectedError
		}

		time.Sleep(2 * time.Second)
	}

	return db, err
}
