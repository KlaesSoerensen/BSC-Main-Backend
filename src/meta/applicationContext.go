package meta

import (
	"otte_main_backend/src/config"
	db "otte_main_backend/src/database"
	"otte_main_backend/src/vitec"
)

type ApplicationContext struct {
	// The application context is a struct that holds all the necessary resources for the application to run.
	ColonyAssetDB            db.ColonyAssetDB
	LanguageDB               db.LanguageDB
	PlayerDB                 db.PlayerDB
	VitecIntegration         *vitec.VitecIntegration
	DDH                      string
	AuthTokenName            string
	MultiplayerServerAddress string
}

func CreateApplicationContext(colonyAssetDB db.ColonyAssetDB, languageDB db.LanguageDB, playerDB db.PlayerDB, vitecIntegration *vitec.VitecIntegration) (*ApplicationContext, error) {
	mpServerAddr, err := config.LoudGet("MULTIPLAYER_SERVER_ADDRESS")
	if err != nil {
		return nil, err
	}

	authTokenName, err := config.LoudGet("AUTH_TOKEN_NAME")
	if err != nil {
		return nil, err
	}

	ddh, err := config.LoudGet("DEFAULT_DEBUG_HEADER")
	if err != nil {
		return nil, err
	}

	return &ApplicationContext{
		ColonyAssetDB:            colonyAssetDB,
		LanguageDB:               languageDB,
		PlayerDB:                 playerDB,
		VitecIntegration:         vitecIntegration,
		DDH:                      ddh,
		AuthTokenName:            authTokenName,
		MultiplayerServerAddress: mpServerAddr,
	}, nil
}
