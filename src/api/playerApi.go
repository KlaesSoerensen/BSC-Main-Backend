package api

import (
	"errors"
	"log"
	"otte_main_backend/src/api/colony"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/database"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/util"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PlayerPreferenceModel represents a single preference item.
type PlayerPreferenceModel struct {
	ID            uint32 `json:"id"`
	Player        uint32 `json:"player" gorm:"column:player;foreignKey:player;references:ID"`
	PreferenceKey string `json:"key" gorm:"column:preferenceKey;foreignKey:preferenceKey;references:preferenceKey"`
	ChosenValue   string `json:"chosenValue" gorm:"column:chosenValue"`
}

func (p *PlayerPreferenceModel) TableName() string {
	return "PlayerPreference"
}

type PlayerPreferenceResponseDTO struct {
	PreferenceKey string `json:"key" gorm:"column:preferenceKey;foreignKey:preferenceKey;references:preferenceKey"`
	ChosenValue   string `json:"chosenValue" gorm:"column:chosenValue"`
}

func (p *PlayerPreferenceResponseDTO) TableName() string {
	return "PlayerPreference"
}

// PlayerPreferencesResponse represents the data returned for a player's preferences.
type PlayerPreferencesResponse struct {
	Preferences []PlayerPreferenceResponseDTO `json:"preferences"`
}

// ColonyInfoResponse represents the data for a single colony.
type ColonyInfoResponse struct {
	ID          uint32              `json:"id"`
	AccLevel    uint32              `json:"accLevel" gorm:"column:accLevel"`
	Name        string              `json:"name"`
	LatestVisit string              `json:"latestVisit" gorm:"column:latestVisit"`
	Assets      []ColonyAssetDTO    `json:"assets" gorm:"foreignKey:ColonyID;references:ID"`
	Locations   []ColonyLocationDTO `json:"locations" gorm:"foreignKey:ColonyID;references:ID"`
}

// ColonyAssetDTO links to TransformDTO with a proper foreign key reference
type ColonyAssetDTO struct {
	ColonyID          uint32       `json:"colonyID" gorm:"column:colony"` // This links back to the colony
	AssetCollectionID uint32       `json:"assetCollectionID" gorm:"column:assetCollection"`
	TransformID       uint32       `json:"-"`                                       // This stores the transform ID
	Transform         TransformDTO `gorm:"foreignKey:TransformID" json:"transform"` // Specify the foreign key
}

// ColonyLocationDTO also links to TransformDTO with a proper foreign key reference
type ColonyLocationDTO struct {
	ColonyID          uint32       `json:"colonyID" gorm:"column:colony"` // This links back to the colony
	Level             uint32       `json:"level"`
	AssetCollectionID uint32       `json:"assetCollectionID" gorm:"column:assetCollection"`
	TransformID       uint32       `json:"-"`                                       // This stores the transform ID
	Transform         TransformDTO `gorm:"foreignKey:TransformID" json:"transform"` // Specify the foreign key
}

// ColonyOverviewResponse represents a collection of colonies.
type ColonyOverviewResponse struct {
	Colonies []ColonyInfoResponse `json:"colonies"`
}

// Apply the Player API routes
func applyPlayerApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Player API] Applying Player API")

	app.Post("/api/v1/player/:playerId/achievement/:achievementId", auth.PrefixOn(appContext, grantPlayerAchievementHandler))

	// Route for fetching a player's preferences by their ID
	app.Get("/api/v1/player/:playerId/preferences", auth.PrefixOn(appContext, getPlayerPreferencesHandler))

	app.Post("/api/v1/player/:playerId/preferences", auth.PrefixOn(appContext, setPlayerPreferenceHandler))

	// Route for fetching colony info by colonyId and playerId
	app.Get("/api/v1/player/:playerId/colony/:colonyId", auth.PrefixOn(appContext, getColonyInfoHandler))

	// Route for fetching overview of all colonies for a player
	app.Get("/api/v1/player/:playerId/colonies", auth.PrefixOn(appContext, getColonyOverviewHandler))

	// Route for creating a new colony
	app.Post("/api/v1/player/:playerId/colony/create", auth.PrefixOn(appContext, createColonyHandler))

	// Route for fetching a single player's info by their ID
	app.Get("/api/v1/player/:playerId", auth.PrefixOn(appContext, getPlayerInfoHandler))
	return nil
}

type SetPreferenceRequestDTO struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AvailablePreferenceModel struct {
	ID              uint32             `gorm:"primaryKey"`
	PreferenceKey   string             `gorm:"column:preferenceKey"`
	AvailableValues util.PGStringArray `gorm:"column:availableValues"`
}

func (a *AvailablePreferenceModel) TableName() string {
	return "AvailablePreference"
}

func setPlayerPreferenceHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerId, parseErr := c.ParamsInt("playerId")
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	var request SetPreferenceRequestDTO
	if err := c.BodyParser(&request); err != nil || request.Key == "" || request.Value == "" {
		c.Response().Header.Set(appContext.DDH, "Invalid request body")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Check if the preference key exists
	var availablePreference AvailablePreferenceModel
	if err := appContext.PlayerDB.Where(`"preferenceKey" = ?`, request.Key).First(&availablePreference).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No such preference")
			return fiber.NewError(fiber.StatusNotFound, "No such preference")
		}

		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	var existingPreference PlayerPreferenceModel
	if err := appContext.PlayerDB.Where(`"player" = ? AND "preferenceKey" = ?`, playerId, request.Key).First(&existingPreference).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Internal server error")
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
	}

	// If the preference already exists, update it
	var updatedPreference = PlayerPreferenceModel{
		ID:            existingPreference.ID,
		Player:        uint32(playerId),
		PreferenceKey: request.Key,
		ChosenValue:   request.Value,
	}

	if setValueErr := appContext.PlayerDB.Table("PlayerPreference").
		Save(&updatedPreference).Error; setValueErr != nil {
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	c.Status(fiber.StatusOK)
	return nil
}

func grantPlayerAchievementHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerId, parseErr := c.ParamsInt("playerId")
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}
	achievementId, parseErr := c.ParamsInt("achievementId")
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid achievement ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid achievement ID "+parseErr.Error())
	}
	var achievement AchievementModel
	if achievementExistErr := appContext.PlayerDB.Where("id = ?", achievementId).First(&achievement).Error; achievementExistErr != nil {
		if !errors.Is(achievementExistErr, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Internal error")
			return fiber.NewError(fiber.StatusNotFound, "Internal error")
		}

		c.Response().Header.Set(appContext.DDH, "Achievement does not exist "+achievementExistErr.Error())
		return fiber.NewError(fiber.StatusNotFound, "Achievement does not exist "+achievementExistErr.Error())
	}

	//Separate check needed as the insert statement doesn't actually error
	//if the player doesn't exist
	if playerExistErr := appContext.PlayerDB.Table("Player").Where("id = ?", playerId).First(&PlayerDTO{}).Error; playerExistErr != nil {
		if !errors.Is(playerExistErr, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Internal error")
			return fiber.NewError(fiber.StatusNotFound, "Internal error")
		}

		c.Response().Header.Set(appContext.DDH, "Player does not exist "+playerExistErr.Error())
		return fiber.NewError(fiber.StatusNotFound, "Player does not exist "+playerExistErr.Error())
	}

	if insertErr := GrantPlayerAchievement(appContext.PlayerDB, uint32(playerId), uint32(achievementId)); insertErr != nil {
		if !errors.Is(insertErr, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Internal error")
			return fiber.NewError(fiber.StatusNotFound, "Internal error")
		}

		c.Response().Header.Set(appContext.DDH, "No such player")
		return fiber.NewError(fiber.StatusInternalServerError, "No such player")
	}

	c.Status(fiber.StatusOK)
	return nil
}

func GrantPlayerAchievement(db database.PlayerDB, playerId uint32, achievementID uint32) error {
	log.Println("[delete me], player id: ", playerId, " achievement id: ", achievementID)
	result := db.Exec(`
        UPDATE "Player"
        SET achievements = ARRAY(SELECT DISTINCT UNNEST(array_append(achievements, ?)))
        WHERE "id" = ?
    `, achievementID, playerId)

	return result.Error
}

func getPlayerPreferencesHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerIdStr := c.Params("playerId")
	playerId, parseErr := strconv.Atoi(playerIdStr)
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID "+parseErr.Error())
	}

	if !doesPlayerExist(appContext.PlayerDB, playerId) {
		c.Response().Header.Set(appContext.DDH, "Player does not exist")
		return fiber.NewError(fiber.StatusNotFound, "Player does not exist")
	}

	// Fetch player preferences and join them with available values
	var preferences []PlayerPreferenceResponseDTO
	if err := appContext.PlayerDB.
		Table(`PlayerPreference`).
		Where(`"PlayerPreference".player = ?`, playerId).
		Select(`"PlayerPreference".id,
				"PlayerPreference"."preferenceKey",
				"PlayerPreference"."chosenValue",
				"AvailablePreference"."availableValues"`).
		Joins(`JOIN "AvailablePreference" ON "PlayerPreference"."preferenceKey" = "AvailablePreference"."preferenceKey"`).
		Find(&preferences).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Preferences not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Preferences not found "+err.Error())
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error")
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Return the preferences in a structured format
	response := PlayerPreferencesResponse{Preferences: preferences}

	c.Status(fiber.StatusOK)
	return c.JSON(response)
}

func doesPlayerExist(db database.PlayerDB, playerId int) bool {
	var player PlayerDTO
	if err := db.Table("Player").Where("id = ?", playerId).First(&player).Error; err != nil {
		return false
	}
	return true
}

func getPlayerInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerId, parseErr := c.ParamsInt("playerId")
	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID: "+strconv.Itoa(playerId))
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID: "+strconv.Itoa(playerId))
	}

	// Fetch player information from the database
	var player PlayerDTO
	if err := appContext.PlayerDB.
		Table("Player").
		Select(`id, "firstName", "lastName", sprite, achievements`).
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

// Model
type ColonyModel struct {
	ID          uint32 `gorm:"primaryKey"`
	Name        string
	AccLevel    uint32 `gorm:"column:accLevel"`
	LatestVisit string `gorm:"column:latestVisit"`
	ColonyCode  uint32 `gorm:"foreignKey:ColonyCode;references:ID;column:colonyCode"`
	Owner       uint32
	Assets      util.PGIntArray
	Locations   util.PGIntArray
}

func (c *ColonyModel) TableName() string {
	return "Colony"
}

type AssetCollectionID struct {
	ID uint32 `gorm:"primaryKey"`
}

func (a *AssetCollectionID) TableName() string {
	return "AssetCollection"
}

type ColonyLocationModel struct {
	ID        uint32 `gorm:"primaryKey"`
	Level     uint32
	Colony    uint32 `gorm:"foreignKey:Colony;references:ID"`
	Transform uint32 `gorm:"foreignKey:Transform;references:ID"`
	Location  uint32 `gorm:"foreignKey:Location;references:ID"`
}

func (a *ColonyLocationModel) TableName() string {
	return "ColonyLocation"
}

// Handler for fetching colony info by playerId and colonyId
func getColonyInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerId, playerIdErr := c.ParamsInt("playerId")
	colonyId, colonyIdErr := c.ParamsInt("colonyId")

	if playerIdErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+playerIdErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID")
	}
	if colonyIdErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+colonyIdErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	// Fetch the colony information
	var colony ColonyModel
	if err := appContext.ColonyAssetDB.
		Where("id = ? AND owner = ?", colonyId, playerId).
		First(&colony).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Colony not found")
		}

		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	type AssetTransformTuple struct {
		Transform         TransformDTO `json:"transform"`
		AssetCollectionID uint32       `json:"assetCollectionID"`
	}
	colonyAssets := make([]AssetTransformTuple, 0, len(colony.Assets))
	for _, colonyAssetID := range colony.Assets {
		var transform TransformDTO
		subQuery := appContext.ColonyAssetDB.
			Table("ColonyAsset").
			Select("transform").
			Where("id = ?", colonyAssetID)
		if err := appContext.ColonyAssetDB.
			Table("Transform").
			Where("id = (?)", subQuery).
			First(&transform).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "No such transform "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "No such transform")
			}

			c.Response().Header.Set(appContext.DDH, "Error fetching transforms "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching transforms")
		}

		var assetCollection AssetCollectionID
		if err := appContext.ColonyAssetDB.
			Table("AssetCollection").
			Select(`id`).
			Where(`id = (SELECT id FROM "ColonyAsset" WHERE id = ?)`, colonyAssetID).
			First(&assetCollection).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "No such AssetCollection "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "No such AssetCollection")
			}

			c.Response().Header.Set(appContext.DDH, "Error fetching assetCollection id's "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching assetCollection id's")
		}

		colonyAssets = append(colonyAssets, AssetTransformTuple{
			Transform:         transform,
			AssetCollectionID: assetCollection.ID,
		})

	}

	type LocationTransformTuple struct {
		Transform  TransformDTO `json:"transform"`
		ID         uint32       `json:"id"`
		LocationID uint32       `json:"locationID" gorm:"foreignKey:Location;references:ID;"`
		Level      uint32       `json:"level"`
	}

	colonyLocations := make([]LocationTransformTuple, 0, len(colony.Locations))
	for _, colonyLocationID := range colony.Locations {
		var colonyLocation ColonyLocationModel
		if err := appContext.ColonyAssetDB.
			Where(`id = ?`, colonyLocationID).
			First(&colonyLocation).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "No such location "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "No such location")
			}

			c.Response().Header.Set(appContext.DDH, "Error fetching location "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching location")
		}

		var transform TransformDTO
		if err := appContext.ColonyAssetDB.
			Where(`id = ?`, colonyLocation.Transform).
			First(&transform).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "No such transform "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "No such transform")
			}

			c.Response().Header.Set(appContext.DDH, "Error fetching transforms "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching transforms")
		}

		colonyLocations = append(colonyLocations, LocationTransformTuple{
			ID:         colonyLocation.ID,
			Transform:  transform,
			LocationID: colonyLocation.Location,
			Level:      colonyLocation.Level,
		})
	}

	toReturn := struct {
		ID          uint32                   `json:"id"`
		AccLevel    uint32                   `json:"accLevel"`
		Name        string                   `json:"name"`
		LatestVisit string                   `json:"latestVisit"`
		Assets      []AssetTransformTuple    `json:"assets"`
		Locations   []LocationTransformTuple `json:"locations"`
	}{
		ID:          colony.ID,
		AccLevel:    colony.AccLevel,
		Name:        colony.Name,
		LatestVisit: colony.LatestVisit,
		Assets:      colonyAssets,
		Locations:   colonyLocations,
	}

	c.Status(fiber.StatusOK)
	return c.JSON(toReturn)
}

// Handler for fetching an overview of all colonies for a player
func getColonyOverviewHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	playerId, playerIdErr := c.ParamsInt("playerId")

	if playerIdErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+playerIdErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID")
	}

	// Fetch all colonies owned by the player, including assets and locations
	var colonies []ColonyModel
	if err := appContext.ColonyAssetDB.
		Where("owner = ?", playerId).
		Find(&colonies).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No colonies found for player "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "No colonies found for player")
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	type ColonyData struct {
		ID          uint32   `json:"id"`
		AccLevel    uint32   `json:"accLevel"`
		Name        string   `json:"name"`
		LatestVisit string   `json:"latestVisit"`
		Assets      []uint32 `json:"assets"`
		Locations   []uint32 `json:"locations"`
	}
	// Prepare the response with the anonymous struct
	var colonyResponses = make([]ColonyData, 0, len(colonies))
	for _, colony := range colonies {
		// Convert util.PGIntArray to []uint32 for assets and locations
		assets := make([]uint32, len(colony.Assets))
		locations := make([]uint32, len(colony.Locations))

		for i, id := range colony.Assets {
			assets[i] = uint32(id)
		}
		for i, id := range colony.Locations {
			locations[i] = uint32(id)
		}

		// Append each formatted colony data
		colonyResponses = append(colonyResponses, ColonyData{
			ID:          colony.ID,
			AccLevel:    colony.AccLevel,
			Name:        colony.Name,
			LatestVisit: colony.LatestVisit,
			Assets:      assets,    // Converted to []uint32
			Locations:   locations, // Converted to []uint32
		})
	}

	// Prepare the colony overview response
	overviewResponse := struct {
		Colonies []ColonyData `json:"colonies"`
	}{
		Colonies: colonyResponses,
	}

	// Return the response
	c.Status(fiber.StatusOK)
	return c.JSON(overviewResponse)
}

// CreateColonyRequest defines the structure for the colony creation request
type CreateColonyRequest struct {
	Name string `json:"name,omitempty"`
}

// Handler for creating a new colony with bare essentials
func createColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	// Get playerId from URL parameters
	playerId, playerIdErr := c.ParamsInt("playerId")
	if playerIdErr != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid player ID "+playerIdErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid player ID")
	}

	// Parse request body for optional colony name
	var request CreateColonyRequest
	if err := c.BodyParser(&request); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Default name if not provided
	colonyName := request.Name
	if colonyName == "" {
		colonyName = "DATA.UNNAMED.COLONY"
	}

	// Start a transaction for the colony insert
	tx := appContext.ColonyAssetDB.Begin()
	if tx.Error != nil {
		c.Response().Header.Set(appContext.DDH, "Error starting transaction: "+tx.Error.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error starting transaction")
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback() // Rollback on panic
		}
	}()

	handleError := func(errMsg string, err error, deleteColony bool, locationIDMap map[uint]uint, transformIDs map[string]uint, newColony *ColonyModel) error {
		tx.Rollback() // Always rollback the current transaction

		// Ensure locations are deleted first, followed by transforms, and then colony itself
		if deleteColony {
			// Extract location IDs
			var locationIDValues []uint
			for _, id := range locationIDMap {
				locationIDValues = append(locationIDValues, id)
			}

			// Delete locations first
			if len(locationIDValues) > 0 {
				appContext.ColonyAssetDB.Where("id IN ?", locationIDValues).Delete(&colony.ColonyLocation{})
			}

			// Extract transform IDs
			var transformIDValues []uint
			for _, id := range transformIDs {
				transformIDValues = append(transformIDValues, id)
			}

			// Delete transforms next
			if len(transformIDValues) > 0 {
				appContext.ColonyAssetDB.Where("id IN ?", transformIDValues).Delete(&colony.Transform{})
			}

			// Finally, delete the colony
			if newColony != nil {
				appContext.ColonyAssetDB.Delete(newColony)
			}
		}

		log.Printf("%s: %v", errMsg, err)
		c.Response().Header.Set(appContext.DDH, errMsg+": "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error creating colony")
	}

	// Create the new colony
	newColony := ColonyModel{
		Name:        colonyName,
		Owner:       uint32(playerId),
		AccLevel:    0,
		LatestVisit: "DATA.UNVISITED.COLONY",
		Assets:      make([]int, 0),
		Locations:   make([]int, 0),
	}

	if err := tx.Create(&newColony).Error; err != nil {
		return handleError("Error creating colony", err, false, nil, nil, nil)
	}

	// Commit colony creation to allow foreign keys in the next steps
	if err := tx.Commit().Error; err != nil {
		return handleError("Error committing colony creation", err, false, nil, nil, &newColony)
	}

	// Start a new transaction for colony locations
	tx = appContext.ColonyAssetDB.Begin()

	// Insert transforms
	transformIDs, err := colony.InsertTransforms(appContext, tx)
	if err != nil {
		return handleError("Error inserting transforms", err, true, nil, transformIDs, &newColony)
	}

	// Insert colony locations and commit them
	locationIDMap, err := colony.InsertColonyLocations(appContext, tx, uint(newColony.ID), transformIDs)
	if err != nil {
		return handleError("Error inserting colony locations", err, true, nil, transformIDs, &newColony)
	}

	// Update the Locations array in newColony
	for _, locationID := range locationIDMap {
		newColony.Locations = append(newColony.Locations, int(locationID))
	}

	// Update the colony record with the new Locations array
	if err := tx.Model(&newColony).Update("Locations", newColony.Locations).Error; err != nil {
		return handleError("Error updating colony locations", err, true, locationIDMap, transformIDs, &newColony)
	}

	// Commit colony locations to ensure they exist for path insertion
	if err := tx.Commit().Error; err != nil {
		return handleError("Error committing colony locations", err, true, locationIDMap, transformIDs, &newColony)
	}

	// Start another transaction for inserting paths
	tx = appContext.ColonyAssetDB.Begin()

	// Insert colony paths using locationIDMap
	if err := colony.InitializeColonyPaths(tx, newColony.ID, locationIDMap); err != nil {
		return handleError("Error initializing colony paths", err, true, locationIDMap, transformIDs, &newColony)
	}

	// Commit the path insertions
	if err := tx.Commit().Error; err != nil {
		return handleError("Error committing colony paths", err, true, locationIDMap, transformIDs, &newColony)
	}

	// Return the newly created colony details
	toReturn := struct {
		ID          uint32 `json:"id"`
		Name        string `json:"name"`
		AccLevel    uint32 `json:"accLevel"`
		LatestVisit string `json:"latestVisit"`
	}{
		ID:          newColony.ID,
		Name:        newColony.Name,
		AccLevel:    newColony.AccLevel,
		LatestVisit: newColony.LatestVisit,
	}

	c.Status(fiber.StatusOK)
	return c.JSON(toReturn)
}
