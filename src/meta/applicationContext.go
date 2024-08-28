package meta

import (
	db "otte_main_backend/src/database"
)

type ApplicationContext struct {
	// The application context is a struct that holds all the necessary resources for the application to run.
	ColonyAssetDB db.ColonyAssetDB
	LanguageDB    db.LanguageDB
	PlayerDB      db.PlayerDB
	DDH           string
}
