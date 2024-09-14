package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PlayerInfoResponse represents the data returned for a player's basic information.
type PlayerInfoResponse struct {
	ID                   uint32          `json:"id"`
	IGN                  string          `json:"IGN"`
	Sprite               uint32          `json:"sprite"`
	Achievements         util.PGIntArray `json:"achievements"`
	HasCompletedTutorial bool            `json:"hasCompletedTutorial"`
}

// PlayerPreference represents a single preference item.
type PlayerPreference struct {
	ID              uint32             `json:"id"`
	PreferenceKey   string             `json:"key" gorm:"column:preferenceKey;foreignKey:preferenceKey;references:preferenceKey"`
	ChosenValue     string             `json:"chosenValue" gorm:"column:chosenValue"`
	AvailableValues util.PGStringArray `json:"availableValues" gorm:"column:availableValues"`
}

// PlayerPreferencesResponse represents the data returned for a player's preferences.
type PlayerPreferencesResponse struct {
	Preferences []PlayerPreference `json:"preferences"`
}

// ColonyInfoResponse represents the data for a single colony.
type ColonyInfoResponse struct {
	ID          uint32              `json:"id"`
	AccLevel    uint32              `json:"accLevel" gorm:"column:accLevel"`
	Name        string              `json:"name"`
	LatestVisit string              `json:"latestVisit" gorm:"column:latestVisit"`
	Assets      []ColonyAssetDTO    `json:"assets" gorm:"foreignKey:ColonyID;references:ID"`    // Define the foreign key for Assets
	Locations   []ColonyLocationDTO `json:"locations" gorm:"foreignKey:ColonyID;references:ID"` // Also specify foreign key for Locations if needed
}

type ColonyAssetDTO struct {
	ColonyID          uint32       `json:"colonyID" gorm:"column:colony"` // This links back to the colony
	AssetCollectionID uint32       `json:"assetCollectionID" gorm:"column:assetCollection"`
	Transform         TransformDTO `json:"transform" gorm:"foreignKey:Transform;references:ID"`
}

type ColonyLocationDTO struct {
	ColonyID          uint32       `json:"colonyID" gorm:"column:colony"` // This links back to the colony
	Level             uint32       `json:"level"`
	AssetCollectionID uint32       `json:"assetCollectionID" gorm:"column:assetCollection"`
	Transform         TransformDTO `json:"transform" gorm:"foreignKey:Transform;references:ID"`
}

// ColonyOverviewResponse represents a collection of colonies.
type ColonyOverviewResponse struct {
	Colonies []ColonyInfoResponse `json:"colonies"`
}

// Apply the Player API routes
func applyPlayerApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Player API] Applying Player API")

	// Route for fetching a single player's info by their ID
	app.Get("/api/v1/player/:playerId", auth.PrefixOn(appContext, getPlayerInfoHandler))

	// Route for fetching a player's preferences by their ID
	app.Get("/api/v1/player/:playerId/preferences", auth.PrefixOn(appContext, getPlayerPreferencesHandler))

	// Route for fetching colony info by colonyId and playerId
	app.Get("/api/v1/player/:playerId/colony/:colonyId", auth.PrefixOn(appContext, getColonyInfoHandler))

	// Route for fetching overview of all colonies for a player
	app.Get("/api/v1/player/:playerId/colonies", auth.PrefixOn(appContext, getColonyOverviewHandler))

	return nil
}

func getPlayerPreferencesHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	// Fetch player preferences and join them with available values
	var preferences []PlayerPreference
	if err := appContext.PlayerDB.
		Table(`PlayerPreference`).
		Where(`"PlayerPreference".player = ?`, playerId).
		Select(`
			"PlayerPreference".id,
			"PlayerPreference"."preferenceKey",
			"PlayerPreference"."chosenValue",
			"AvailablePreference"."availableValues"`).
		Joins(`
			JOIN "AvailablePreference" ON "PlayerPreference"."preferenceKey" = "AvailablePreference"."preferenceKey"`).
		Find(&preferences).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Preferences not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Preferences not found "+err.Error())
		}
		//Gorm exposes secrets in err when DB is down, so it can't be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Return the preferences in a structured format
	response := PlayerPreferencesResponse{Preferences: preferences}

	c.Status(fiber.StatusOK)
	return c.JSON(response)
}

func getPlayerInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	// Fetch player information from the database
	var player PlayerInfoResponse
	if err := appContext.PlayerDB.
		Table("Player").
		Select(`id, "IGN", sprite, achievements`).
		Where("id = ?", playerId).
		First(&player).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Player not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Player not found "+err.Error())
		}

		// Gorm exposes secrets in err when DB is down, so it can't be included in the response
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Compute the 'HasCompletedTutorial' field
	// Achievement id 1 is always the tutorial achievement
	player.HasCompletedTutorial = util.ArrayContains(player.Achievements, 1)

	c.Status(fiber.StatusOK)
	return c.JSON(player)
}

// Handler for fetching colony info by playerId and colonyId
func getColonyInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	colonyIdStr := c.Params("colonyId")

	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID")
	}
	colonyId, parseErr := strconv.Atoi(colonyIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	// Fetch the colony information
	var colony ColonyInfoResponse
	if err := appContext.ColonyAssetDB.
		Table("Colony").
		Select(`id, "accLevel", name, "latestVisit"`).
		Where("id = ? AND owner = ?", colonyId, playerId).
		First(&colony).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Colony not found")
		}

		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Fetch assets for the colony, including the TransformDTO automatically
	var assets []ColonyAssetDTO
	if err := appContext.ColonyAssetDB.
		Preload("Transform"). // Preload the TransformDTO
		Table("ColonyAsset").
		Select(`assetCollection`).
		Where("colony = ?", colonyId).
		Find(&assets).Error; err != nil {
		c.Response().Header.Set(appContext.DDH, "Error fetching assets "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error fetching assets")
	}

	// Fetch locations for the colony, including the TransformDTO automatically
	var locations []ColonyLocationDTO
	if err := appContext.ColonyAssetDB.
		Preload("Transform"). // Preload the TransformDTO
		Table("ColonyLocation").
		Select(`level, location`).
		Where("colony = ?", colonyId).
		Find(&locations).Error; err != nil {
		c.Response().Header.Set(appContext.DDH, "Error fetching locations "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error fetching locations")
	}

	// Attach assets and locations to the colony response
	colony.Assets = assets
	colony.Locations = locations

	// Return the response in the desired format
	return c.JSON(colony)
}

// Handler for fetching all colonies for a specific player
func getColonyOverviewHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID")
	}

	// Fetch all colonies for the player
	var colonies []ColonyInfoResponse
	if err := appContext.ColonyAssetDB.
		Table("Colony").
		Where("owner = ?", playerId).
		Find(&colonies).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No colonies found for player "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "No colonies found for player")
		}

		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Fetch assets and locations for each colony
	for i := range colonies {
		colonyId := colonies[i].ID

		// Fetch assets for this colony, including the TransformDTO automatically
		var assets []ColonyAssetDTO
		if err := appContext.ColonyAssetDB.
			Preload("Transform"). // Preload the TransformDTO
			Table("ColonyAsset").
			Select(`assetCollection`).
			Where("colony = ?", colonyId).
			Find(&assets).Error; err != nil {
			c.Response().Header.Set(appContext.DDH, "Error fetching assets for colony "+strconv.Itoa(int(colonyId))+" "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching assets for colony")
		}
		colonies[i].Assets = assets

		// Fetch locations for this colony, including the TransformDTO automatically
		var locations []ColonyLocationDTO
		if err := appContext.ColonyAssetDB.
			Preload("Transform"). // Preload the TransformDTO
			Table("ColonyLocation").
			Select(`"level", location`).
			Where("colony = ?", colonyId).
			Find(&locations).Error; err != nil {
			c.Response().Header.Set(appContext.DDH, "Error fetching locations for colony "+strconv.Itoa(int(colonyId))+" "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching locations for colony")
		}
		colonies[i].Locations = locations
	}

	c.Status(fiber.StatusOK)
	return c.JSON(ColonyOverviewResponse{Colonies: colonies})
}
