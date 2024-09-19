package api

import (
	"errors"
	"log"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type LocationAppearanceDTO struct {
	Level             uint32 `json:"level"`
	AssetCollectionID uint32 `json:"assetCollectionID" gorm:"column:assetCollection"`
	Assets            []struct {
		Transform TransformDTO              `json:"transform"`
		Asset     MinimizedAssetLocationDTO `json:"asset"`
	} `json:"assets"`
}

// LocationInfoResponse represents the data for a location
type LocationInfoResponse struct {
	ID          uint32                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Appearances []LocationAppearanceDTO `json:"appearances"`
	MinigameID  uint32                  `json:"minigameID"`
}

// Apply the Location API routes
func applyLocationApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Location API] Applying location API")

	// Route for fetching location info by locationID
	app.Get("/api/v1/location/:locationID", getLocationInfoHandler(appContext))

	// Route for fetching full location info by locationID
	app.Get("/api/v1/location/:locationID/full", getLocationFullInfoHandler(appContext)) // <-- New Route

	return nil
}

// Handler for fetching location info by locationID
func getLocationInfoHandler(appContext *meta.ApplicationContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get locationID from the URL
		locationID, locationIDErr := c.ParamsInt("locationID")
		if locationIDErr != nil {
			c.Response().Header.Set(appContext.DDH, "Invalid location ID "+locationIDErr.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
		}

		// Fetch the location information
		var location LocationModel
		if err := appContext.ColonyAssetDB.
			Where("id = ?", locationID).
			First(&location).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "Location not found "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "Location not found")
			}

			c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		// Fetch the location appearances
		var appearances []LocationAppearanceDTO
		if err := appContext.ColonyAssetDB.
			Table("LocationAppearance").
			Where("location = ?", locationID).
			Find(&appearances).Error; err != nil {
			c.Response().Header.Set(appContext.DDH, "Error fetching location appearances "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Error fetching location appearances")
		}

		// Prepare the response
		toReturn := LocationInfoResponse{
			ID:          location.ID,
			Name:        location.Name,
			Description: location.Description,
			Appearances: appearances,
			MinigameID:  location.MinigameID, // Assuming minigameID is stored in LocationModel
		}

		c.Status(fiber.StatusOK)
		return c.JSON(toReturn)
	}
}

// LocationModel represents the Location table
type LocationModel struct {
	ID          uint32 `gorm:"primaryKey"`
	Name        string
	Description string
	MinigameID  uint32                    `gorm:"column:minigame"`
	Minigame    MinigameModel             `gorm:"foreignKey:MinigameID;references:ID"`
	Appearances []LocationAppearanceModel `gorm:"foreignKey:LocationID;references:ID"`
}

func (l *LocationModel) TableName() string {
	return "Location"
}

// MinigameDTO represents detailed minigame information
type MinigameDTO struct {
	ID           uint32                          `json:"id"`
	Name         string                          `json:"name"`
	Description  string                          `json:"description"`
	IconID       uint32                          `json:"iconID"`
	Difficulties []MinigameDifficultyLocationDTO `json:"difficulties"`
}

// LocationFullInfoResponse represents the full location data
type LocationFullInfoResponse struct {
	ID          uint32                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Appearances []LocationAppearanceDTO `json:"appearances"`
	Minigame    MinigameDTO             `json:"minigame"`
}

type MinigameModel struct {
	ID           uint32 `gorm:"primaryKey"`
	Name         string
	Description  string
	IconID       uint32                    `gorm:"foreignKey:IconID"`
	Icon         GraphicalAssetModel       `gorm:"foreignKey:IconID;references:ID"`
	Difficulties []MinigameDifficultyModel `gorm:"foreignKey:MinigameID;references:ID"` // Define the foreign key
}

func (MinigameModel) TableName() string {
	return "MiniGame" // Matches the database table name
}

type MinigameDifficultyModel struct {
	ID          uint32 `gorm:"primaryKey"`
	Name        string
	Description string
	IconID      uint32              `gorm:"foreignKey:IconID"`
	Icon        GraphicalAssetModel `gorm:"foreignKey:IconID;references:ID"`
	MinigameID  uint32              `gorm:"column:minigame"` // This is the foreign key field
}

func (MinigameDifficultyModel) TableName() string {
	return "MiniGameDifficulty" // Matches the database table name
}

// MinigameDifficultyLocationDTO represents minigame difficulties
type MinigameDifficultyLocationDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconID      uint32 `json:"iconID"`
}

// GraphicalAssetModel represents the GraphicalAsset table
type GraphicalAssetModel struct {
	ID      uint32 `gorm:"primaryKey"`
	Alias   string
	Type    string
	UseCase string `gorm:"column:useCase"`
	Width   int
	Height  int
}

func (GraphicalAssetModel) TableName() string {
	return "GraphicalAsset" // Matches the database table name
}

// LocationAppearanceModel represents the LocationAppearance table
type LocationAppearanceModel struct {
	ID                uint32 `gorm:"primaryKey"`
	Level             uint32
	LocationID        uint32               `gorm:"column:location"`
	AssetCollectionID uint32               `gorm:"column:assetCollection"`
	AssetCollection   AssetCollectionModel `gorm:"foreignKey:AssetCollectionID;references:ID"`
}

func (LocationAppearanceModel) TableName() string {
	return "LocationAppearance" // Matches the database table name
}

type AssetCollectionModel struct {
	ID                uint32 `gorm:"primaryKey"`
	Name              string
	UseCase           string                 `gorm:"column:useCase"`
	CollectionEntries []CollectionEntryModel `gorm:"foreignKey:AssetCollectionID;references:ID"`
}

func (AssetCollectionModel) TableName() string {
	return "AssetCollection" // Matches the database table name
}

type CollectionEntryModel struct {
	ID                uint32              `gorm:"primaryKey"`
	GraphicalAssetID  uint32              `gorm:"column:graphicalAsset"`
	GraphicalAsset    GraphicalAssetModel `gorm:"foreignKey:GraphicalAssetID;references:ID"`
	TransformID       uint32              `gorm:"column:transform"`
	Transform         TransformModel      `gorm:"foreignKey:TransformID;references:ID"`
	AssetCollectionID uint32              `gorm:"column:assetCollection"`
}

func (CollectionEntryModel) TableName() string {
	return "CollectionEntry" // Matches the database table name
}

// TransformModel represents the Transform table
type TransformModel struct {
	ID      uint32  `gorm:"primaryKey;column:id" json:"id"`
	XScale  float32 `gorm:"column:xScale" json:"xScale"`
	YScale  float32 `gorm:"column:yScale" json:"yScale"`
	XOffset float32 `gorm:"column:xOffset" json:"xOffset"`
	YOffset float32 `gorm:"column:yOffset" json:"yOffset"`
	ZIndex  int     `gorm:"column:zIndex" json:"zIndex"`
}

func (TransformModel) TableName() string {
	return "Transform" // Matches the database table name
}

// MinimizedAssetLocationDTO represents an asset with minimal information
type MinimizedAssetLocationDTO struct {
	ID     uint32 `json:"id"`
	Alias  string `json:"alias"`
	Type   string `json:"type"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func getLocationFullInfoHandler(appContext *meta.ApplicationContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get locationID from the URL
		locationID, locationIDErr := c.ParamsInt("locationID")
		if locationIDErr != nil {
			c.Response().Header.Set(appContext.DDH, "Invalid location ID "+locationIDErr.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
		}

		// Fetch the location information, including minigame and appearances
		var location LocationModel
		if err := appContext.ColonyAssetDB.
			Preload("Minigame.Icon").
			Preload("Minigame.Difficulties.Icon").
			Preload("Appearances.AssetCollection.CollectionEntries.GraphicalAsset").
			Preload("Appearances.AssetCollection.CollectionEntries.Transform").
			Where("id = ?", locationID).
			First(&location).Error; err != nil {

			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "Location not found "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "Location not found")
			}

			c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		// Prepare the appearances response
		appearances := make([]LocationAppearanceDTO, len(location.Appearances))
		for i, appearance := range location.Appearances {
			assets := make([]struct {
				Transform TransformDTO              `json:"transform"`
				Asset     MinimizedAssetLocationDTO `json:"asset"`
			}, len(appearance.AssetCollection.CollectionEntries))

			for j, entry := range appearance.AssetCollection.CollectionEntries {
				assets[j] = struct {
					Transform TransformDTO              `json:"transform"`
					Asset     MinimizedAssetLocationDTO `json:"asset"`
				}{
					Transform: TransformDTO{
						XScale:  entry.Transform.XScale,
						YScale:  entry.Transform.YScale,
						XOffset: entry.Transform.XOffset,
						YOffset: entry.Transform.YOffset,
						ZIndex:  uint32(entry.Transform.ZIndex),
					},
					Asset: MinimizedAssetLocationDTO{
						ID:     entry.GraphicalAsset.ID,
						Alias:  entry.GraphicalAsset.Alias,
						Type:   entry.GraphicalAsset.Type,
						Width:  entry.GraphicalAsset.Width,
						Height: entry.GraphicalAsset.Height,
					},
				}
			}

			appearances[i] = LocationAppearanceDTO{
				Level:  appearance.Level,
				Assets: assets,
			}
		}

		// Prepare the minigame response
		minigame := MinigameDTO{
			ID:           location.Minigame.ID,
			Name:         location.Minigame.Name,
			Description:  location.Minigame.Description,
			IconID:       location.Minigame.Icon.ID,
			Difficulties: []MinigameDifficultyLocationDTO{},
		}

		// Handle difficulties
		for _, difficulty := range location.Minigame.Difficulties {
			minigame.Difficulties = append(minigame.Difficulties, MinigameDifficultyLocationDTO{
				Name:        difficulty.Name,
				Description: difficulty.Description,
				IconID:      difficulty.Icon.ID,
			})
		}

		// Prepare the full location response
		toReturn := LocationFullInfoResponse{
			ID:          location.ID,
			Name:        location.Name,
			Description: location.Description,
			Appearances: appearances,
			Minigame:    minigame,
		}

		c.Status(fiber.StatusOK)
		return c.JSON(toReturn)
	}
}
