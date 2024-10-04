package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func applyLocationApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Location API] Applying location API")

	// Route for fetching basic location info by locationID
	app.Get("/api/v1/location/:locationID", auth.PrefixOn(appContext, getLocationInfoHandler))

	// Route for fetching full location info by locationID
	app.Get("/api/v1/location/:locationID/full", auth.PrefixOn(appContext, getLocationFullInfoHandler))

	return nil
}

// LocationFullInfoResponse represents the full location data
type LocationFullInfoResponse struct {
	ID          uint32                  `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Appearances []LocationAppearanceDTO `json:"appearances"`
	Minigame    MinigameDTO             `json:"minigame"`
}

type LocationAppearanceDTO struct {
	Level     uint32 `json:"level"`
	SplashArt uint32 `json:"splashArt" gorm:"column:splashArt"`
	Assets    []struct {
		Transform TransformLocationDTO      `json:"transform"`
		Asset     MinimizedAssetLocationDTO `json:"asset"`
	} `json:"assets"`
}

type MinimizedAssetLocationDTO struct {
	ID     uint32        `json:"id"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
	Alias  string        `json:"alias"`
	Type   string        `json:"type"`
	LODs   []AssetLODDTO `json:"LODs"`
}

type AssetLODDTO struct {
	DetailLevel uint32 `json:"detailLevel"`
	ID          uint32 `json:"id"`
}

type TransformLocationDTO struct {
	XScale  float32 `json:"xScale"`
	YScale  float32 `json:"yScale"`
	XOffset float32 `json:"xOffset"`
	YOffset float32 `json:"yOffset"`
	ZIndex  uint32  `json:"zIndex"`
}

type MinigameDTO struct {
	ID           uint32                          `json:"id"`
	Name         string                          `json:"name"`
	Description  string                          `json:"description"`
	IconID       uint32                          `json:"iconID"`
	Difficulties []MinigameDifficultyLocationDTO `json:"difficulties"`
}

type MinigameDifficultyLocationDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IconID      uint32 `json:"iconID"`
}

// LocationModel represents the Location table
type LocationModel struct {
	ID          uint32                    `gorm:"column:id;primaryKey"`
	Name        string                    `gorm:"column:name"`
	Description string                    `gorm:"column:description"`
	MinigameID  uint32                    `gorm:"column:minigame"`
	Minigame    MinigameModel             `gorm:"foreignKey:MinigameID"`
	Appearances []LocationAppearanceModel `gorm:"foreignKey:LocationID"`
}

func (lm *LocationModel) TableName() string {
	return "Location"
}

type LocationAppearanceModel struct {
	ID                uint32               `gorm:"column:id;primaryKey"`
	Level             int                  `gorm:"column:level"`
	LocationID        uint32               `gorm:"column:location"`
	SplashArt         uint32               `json:"splashArt" gorm:"column:splashArt"`
	AssetCollectionID uint32               `gorm:"column:assetCollection"`
	AssetCollection   AssetCollectionModel `gorm:"foreignKey:AssetCollectionID"`
}

func (lam *LocationAppearanceModel) TableName() string {
	return "LocationAppearance"
}

type AssetCollectionModel struct {
	ID                uint32                 `gorm:"column:id;primaryKey"`
	Name              string                 `gorm:"column:name"`
	UseCase           string                 `gorm:"column:useCase"`
	CollectionEntries []CollectionEntryModel `gorm:"foreignKey:AssetCollectionID"`
}

func (acm *AssetCollectionModel) TableName() string {
	return "AssetCollection"
}

type CollectionEntryModel struct {
	ID                uint32              `gorm:"column:id;primaryKey"`
	GraphicalAssetID  uint32              `gorm:"column:graphicalAsset"`
	GraphicalAsset    GraphicalAssetModel `gorm:"foreignKey:GraphicalAssetID"`
	TransformID       uint32              `gorm:"column:transform"`
	Transform         TransformModel      `gorm:"foreignKey:TransformID"`
	AssetCollectionID uint32              `gorm:"column:assetCollection"`
}

func (cem *CollectionEntryModel) TableName() string {
	return "CollectionEntry"
}

type GraphicalAssetModel struct {
	ID      uint32     `gorm:"column:id;primaryKey"`
	Alias   string     `gorm:"column:alias"`
	Type    string     `gorm:"column:type"`
	UseCase string     `gorm:"column:useCase"`
	Width   int        `gorm:"column:width"`
	Height  int        `gorm:"column:height"`
	LODs    []LODModel `gorm:"foreignKey:GraphicalAssetID"`
}

func (ga *GraphicalAssetModel) TableName() string {
	return "GraphicalAsset"
}

type LODModel struct {
	ID               uint32 `gorm:"column:id;primaryKey"`
	DetailLevel      int    `gorm:"column:detailLevel"`
	Type             string `gorm:"column:type"`
	GraphicalAssetID uint32 `gorm:"column:graphicalAsset"`
}

func (lm *LODModel) TableName() string {
	return "LOD"
}

type MinigameModel struct {
	ID           uint32                    `gorm:"column:id;primaryKey"`
	Name         string                    `gorm:"column:name"`
	Description  string                    `gorm:"column:description"`
	IconID       uint32                    `gorm:"column:icon"`
	Icon         GraphicalAssetModel       `gorm:"foreignKey:IconID"`
	Difficulties []MinigameDifficultyModel `gorm:"foreignKey:MinigameID"`
}

func (mm *MinigameModel) TableName() string {
	return "MiniGame"
}

type MinigameDifficultyModel struct {
	ID          uint32              `gorm:"column:id;primaryKey"`
	Name        string              `gorm:"column:name"`
	Description string              `gorm:"column:description"`
	IconID      uint32              `gorm:"column:icon"`
	Icon        GraphicalAssetModel `gorm:"foreignKey:IconID"`
	MinigameID  uint32              `gorm:"column:minigame"`
}

func (mdm *MinigameDifficultyModel) TableName() string {
	return "MiniGameDifficulty"
}

func getLocationInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	locationID, err := c.ParamsInt("locationID")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid location ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
	}

	var location LocationModel
	if err := appContext.ColonyAssetDB.
		Preload("Appearances").
		Where("id = ?", locationID).
		First(&location).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Location not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Location not found")
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	type LocationAppearanceDTO struct {
		Level             uint32 `json:"level"`
		SplashArt         uint32 `json:"splashArt" gorm:"column:splashArt"`
		AssetCollectionID uint32 `json:"assetCollectionID"`
	}

	appearances := make([]LocationAppearanceDTO, len(location.Appearances))

	for i, appearance := range location.Appearances {
		appearances[i] = LocationAppearanceDTO{
			Level:             uint32(appearance.Level),
			AssetCollectionID: appearance.AssetCollectionID,
			SplashArt:         appearance.SplashArt,
		}
	}

	response := struct {
		ID          uint32                  `json:"id"`
		Name        string                  `json:"name"`
		Description string                  `json:"description"`
		Appearances []LocationAppearanceDTO `json:"appearances"`
		MinigameID  uint32                  `json:"minigameID"`
	}{
		ID:          location.ID,
		Name:        location.Name,
		Description: location.Description,
		Appearances: appearances,
		MinigameID:  location.MinigameID,
	}

	return c.JSON(response)
}

// Handler for full location info
func getLocationFullInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {

	locationID, err := c.ParamsInt("locationID")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid location ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid location ID")
	}

	var location LocationModel
	if err := appContext.ColonyAssetDB.
		Preload("Minigame.Icon").
		Preload("Minigame.Difficulties.Icon").
		Preload("Appearances.AssetCollection.CollectionEntries.GraphicalAsset.LODs").
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

	appearances := make([]LocationAppearanceDTO, len(location.Appearances))
	for i, appearance := range location.Appearances {
		assets := make([]struct {
			Transform TransformLocationDTO      `json:"transform"`
			Asset     MinimizedAssetLocationDTO `json:"asset"`
		}, len(appearance.AssetCollection.CollectionEntries))

		for j, entry := range appearance.AssetCollection.CollectionEntries {
			lods := make([]AssetLODDTO, len(entry.GraphicalAsset.LODs))
			for k, lod := range entry.GraphicalAsset.LODs {
				lods[k] = AssetLODDTO{
					DetailLevel: uint32(lod.DetailLevel),
					ID:          lod.ID,
				}
			}

			assets[j] = struct {
				Transform TransformLocationDTO      `json:"transform"`
				Asset     MinimizedAssetLocationDTO `json:"asset"`
			}{
				Transform: TransformLocationDTO{
					XScale:  entry.Transform.XScale,
					YScale:  entry.Transform.YScale,
					XOffset: entry.Transform.XOffset,
					YOffset: entry.Transform.YOffset,
					ZIndex:  uint32(entry.Transform.ZIndex),
				},
				Asset: MinimizedAssetLocationDTO{
					ID:     entry.GraphicalAsset.ID,
					Width:  entry.GraphicalAsset.Width,
					Height: entry.GraphicalAsset.Height,
					Alias:  entry.GraphicalAsset.Alias,
					Type:   entry.GraphicalAsset.Type,
					LODs:   lods,
				},
			}
		}

		appearances[i] = LocationAppearanceDTO{
			Level:  uint32(appearance.Level),
			Assets: assets,
		}
	}

	minigame := MinigameDTO{
		ID:           location.Minigame.ID,
		Name:         location.Minigame.Name,
		Description:  location.Minigame.Description,
		IconID:       location.Minigame.IconID,
		Difficulties: []MinigameDifficultyLocationDTO{},
	}

	for _, difficulty := range location.Minigame.Difficulties {
		minigame.Difficulties = append(minigame.Difficulties, MinigameDifficultyLocationDTO{
			Name:        difficulty.Name,
			Description: difficulty.Description,
			IconID:      difficulty.IconID,
		})
	}

	response := LocationFullInfoResponse{
		ID:          location.ID,
		Name:        location.Name,
		Description: location.Description,
		Appearances: appearances,
		Minigame:    minigame,
	}

	return c.JSON(response)
}
