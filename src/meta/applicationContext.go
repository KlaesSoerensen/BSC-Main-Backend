package meta

import (
	"fmt"
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
	// Address to use when accessed by main backend
	InternalMultiplayerServerAddress string
	// Address to use when accessed by anyone else
	ExternalMultiplayerServerAddress string
}

// Returns: internal, external, error
func getMPAddresses() (string, string, error) {
	internalMPAddrHost, err := config.LoudGet("MULTIPLAYER_BACKEND_HOST_INTERNAL")
	if err != nil {
		return "", "", err
	}
	externalMPAddrHost, err := config.LoudGet("MULTIPLAYER_BACKEND_HOST_EXTERNAL")
	if err != nil {
		return "", "", err
	}
	internalMPAddrPort, err := config.GetInt("MULTIPLAYER_BACKEND_PORT_INTERNAL")
	if err != nil {
		return "", "", err
	}
	externalMPAddrPort, err := config.GetInt("MULTIPLAYER_BACKEND_PORT_EXTERNAL")
	if err != nil {
		return "", "", err
	}
	internalMPAddr := fmt.Sprintf("http://%s:%d", internalMPAddrHost, internalMPAddrPort)
	externalMPAddr := fmt.Sprintf("http://%s:%d", externalMPAddrHost, externalMPAddrPort)
	return internalMPAddr, externalMPAddr, nil
}

func CreateApplicationContext(colonyAssetDB db.ColonyAssetDB, languageDB db.LanguageDB, playerDB db.PlayerDB, vitecIntegration *vitec.VitecIntegration) (*ApplicationContext, error) {
	internalMPAddr, externalMPAddr, err := getMPAddresses()
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
		ColonyAssetDB:                    colonyAssetDB,
		LanguageDB:                       languageDB,
		PlayerDB:                         playerDB,
		VitecIntegration:                 vitecIntegration,
		DDH:                              ddh,
		AuthTokenName:                    authTokenName,
		InternalMultiplayerServerAddress: internalMPAddr,
		ExternalMultiplayerServerAddress: externalMPAddr,
	}, nil
}
