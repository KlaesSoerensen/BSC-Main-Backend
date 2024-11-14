package colony

import (
	"fmt"
	"math"
	"math/rand"

	"gorm.io/gorm"
)

type GraphicalAsset struct {
	ID      uint32 `gorm:"column:id;primaryKey"`
	Alias   string `gorm:"column:alias"`
	Type    string `gorm:"column:type"`
	UseCase string `gorm:"column:useCase"`
	Width   int    `gorm:"column:width"`
	Height  int    `gorm:"column:height"`
}

type ColonyAssetInsertDTO struct {
	ID                uint32 `gorm:"column:id;primaryKey"`
	AssetCollectionID uint32 `json:"assetCollection" gorm:"column:assetCollection"`
	Transform         uint   `json:"transform" gorm:"column:transform"`
	Colony            uint32 `json:"colony" gorm:"column:colony"`
}

func (GraphicalAsset) TableName() string {
	return "GraphicalAsset"
}

func (cai *ColonyAssetInsertDTO) TableName() string {
	return "ColonyAsset"
}

const (
	BaseGroundTileCollectionID = 10001 // The ID of the AssetCollection containing the ground tile
	DecorationCollectionMinID  = 10002 // Start of decoration AssetCollection IDs
	DecorationCollectionMaxID  = 10026 // End of decoration AssetCollection IDs (inclusive)
	paddingMultiplier          = 1.5   // Increase coverage beyond the structure bounds
)

// Helper function to get random decoration collection ID
func getRandomDecorationCollectionID() uint32 {
	return uint32(DecorationCollectionMinID + rand.Intn(DecorationCollectionMaxID-DecorationCollectionMinID+1))
}

func createTileTransform(x, y float64, isDecoration bool) Transform {
	zIndex := 0
	if isDecoration {
		zIndex = 1
	}
	return Transform{
		XScale:  1,
		YScale:  1,
		XOffset: x,
		YOffset: y,
		ZIndex:  zIndex,
	}
}

func InsertColonyAssets(tx *gorm.DB, colonyID uint32, boundingBox *BoundingBox) ([]int, error) {
	var insertedAssetIDs []int

	// Get the base tile asset information
	var baseTile GraphicalAsset
	if err := tx.Where("id = ?", 8001).First(&baseTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching base tile asset: %w", err)
	}

	// Calculate the expanded area with padding
	// BoundingBox already includes the topLevelDistanceScalar
	expandedMinX := boundingBox.MinX - (math.Abs(boundingBox.MaxX-boundingBox.MinX) * paddingMultiplier)
	expandedMinY := boundingBox.MinY - (math.Abs(boundingBox.MaxY-boundingBox.MinY) * paddingMultiplier)
	expandedMaxX := boundingBox.MaxX + (math.Abs(boundingBox.MaxX-boundingBox.MinX) * paddingMultiplier)
	expandedMaxY := boundingBox.MaxY + (math.Abs(boundingBox.MaxY-boundingBox.MinY) * paddingMultiplier)

	// Calculate tile dimensions with overlap
	tileWidth := float64(baseTile.Width) * 0.95
	tileHeight := float64(baseTile.Height) * 0.95

	// Calculate number of tiles needed for the expanded area
	tilesX := math.Ceil((expandedMaxX - expandedMinX) / (tileWidth * 0.95))
	tilesY := math.Ceil((expandedMaxY - expandedMinY) / (tileHeight * 0.95))

	// Adjust starting position to center the grid
	offsetX := ((tilesX * tileWidth * 0.95) - (expandedMaxX - expandedMinX)) / 2
	offsetY := ((tilesY * tileHeight * 0.95) - (expandedMaxY - expandedMinY)) / 2
	startX := expandedMinX - offsetX
	startY := expandedMinY - offsetY

	// Create and insert colony assets
	for i := 0; i < int(tilesX); i++ {
		for j := 0; j < int(tilesY); j++ {
			// Calculate tile position with overlap
			tileX := startX + (float64(i) * tileWidth * 0.95)
			tileY := startY + (float64(j) * tileHeight * 0.95)

			// Create base tile transform
			tileTransform := createTileTransform(tileX, tileY, false)
			if err := tx.Create(&tileTransform).Error; err != nil {
				return nil, fmt.Errorf("error creating tile transform: %w", err)
			}

			// Create base tile asset
			baseTileAsset := ColonyAssetInsertDTO{
				AssetCollectionID: BaseGroundTileCollectionID,
				Transform:         tileTransform.ID,
				Colony:            colonyID,
			}

			if err := tx.Create(&baseTileAsset).Error; err != nil {
				return nil, fmt.Errorf("error creating base tile asset: %w", err)
			}
			insertedAssetIDs = append(insertedAssetIDs, int(baseTileAsset.ID))

			// Add 0-3 decorations per tile
			numDecorations := rand.Intn(4)
			for d := 0; d < numDecorations; d++ {
				decorationCollectionID := getRandomDecorationCollectionID()

				// Calculate decoration offsets relative to tile
				offsetX := rand.Float64() * (tileWidth * 0.8)  // 80% of tile width
				offsetY := rand.Float64() * (tileHeight * 0.8) // 80% of tile height

				// Create decoration transform
				decorTransform := createTileTransform(tileX+offsetX, tileY+offsetY, true)
				if err := tx.Create(&decorTransform).Error; err != nil {
					return nil, fmt.Errorf("error creating decoration transform: %w", err)
				}

				// Create decoration asset
				decorationAsset := ColonyAssetInsertDTO{
					AssetCollectionID: decorationCollectionID,
					Transform:         decorTransform.ID,
					Colony:            colonyID,
				}

				if err := tx.Create(&decorationAsset).Error; err != nil {
					return nil, fmt.Errorf("error creating decoration asset: %w", err)
				}
				insertedAssetIDs = append(insertedAssetIDs, int(decorationAsset.ID))
			}
		}
	}

	return insertedAssetIDs, nil
}
