package meta

import (
	"otte_main_backend/src/config"
	db "otte_main_backend/src/database"
	"otte_main_backend/src/vitec"
)

type ApplicationContext struct {
	// The application context is a struct that holds all the necessary resources for the application to run.
	ColonyAssetDB    db.ColonyAssetDB
	LanguageDB       db.LanguageDB
	PlayerDB         db.PlayerDB
	VitecIntegration *vitec.VitecIntegration
	DDH              string
	AuthTokenName    string
}

func CreateApplicationContext(colonyAssetDB db.ColonyAssetDB, languageDB db.LanguageDB, playerDB db.PlayerDB, vitecIntegration *vitec.VitecIntegration, ddh string) *ApplicationContext {
	return &ApplicationContext{
		ColonyAssetDB:    colonyAssetDB,
		LanguageDB:       languageDB,
		PlayerDB:         playerDB,
		VitecIntegration: vitecIntegration,
		DDH:              ddh,
		AuthTokenName:    config.GetOr("AUTH_TOKEN_NAME", "OTTE-Token"),
	}
}
