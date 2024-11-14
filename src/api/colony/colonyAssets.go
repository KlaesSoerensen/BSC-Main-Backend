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
	BaseGroundTileCollectionID = 10027 // The ID of the AssetCollection containing the ground tile
	DecorationCollectionMinID  = 10002 // Start of decoration AssetCollection IDs
	DecorationCollectionMaxID  = 10026 // End of decoration AssetCollection IDs (inclusive)
	horizontalPadding          = 0.5   // 50% horizontal padding
	verticalPadding            = 0.25  // 25% vertical padding
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
	// Get the base tile asset information
	var baseTile GraphicalAsset
	if err := tx.Where("id = ?", 8027).First(&baseTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching base tile asset: %w", err)
	}

	// Calculate tile dimensions with overlap
	tileWidth := float64(baseTile.Width) * 0.95
	tileHeight := float64(baseTile.Height) * 0.95

	// Calculate the expanded area with different padding for horizontal and vertical
	// Horizontal extends both ways
	expandedMinX := boundingBox.MinX - (math.Abs(boundingBox.MaxX-boundingBox.MinX) * horizontalPadding)
	expandedMaxX := boundingBox.MaxX + (math.Abs(boundingBox.MaxX-boundingBox.MinX) * horizontalPadding)

	// Vertical: Start well below top boundary and extend downward
	expandedMinY := boundingBox.MinY + tileHeight*0.875
	expandedMaxY := boundingBox.MaxY + (math.Abs(boundingBox.MaxY-boundingBox.MinY) * verticalPadding * 2)

	// Calculate number of tiles needed for the expanded area
	tilesX := math.Ceil((expandedMaxX - expandedMinX) / (tileWidth * 0.95))
	tilesY := math.Ceil((expandedMaxY - expandedMinY) / (tileWidth * 0.95))

	// Adjust starting position - center horizontally and start below top boundary
	offsetX := ((tilesX * tileWidth * 0.95) - (expandedMaxX - expandedMinX)) / 2
	startX := expandedMinX - offsetX
	startY := expandedMinY // Start below top boundary

	// Prepare batch arrays
	var transforms []Transform
	var assets []ColonyAssetInsertDTO
	var insertedAssetIDs []int

	// Pre-allocate slices with approximate size
	estimatedSize := int(tilesX * tilesY * 2) // Base tiles + estimated decorations
	transforms = make([]Transform, 0, estimatedSize)
	assets = make([]ColonyAssetInsertDTO, 0, estimatedSize)

	// Create all tile transforms and assets first
	for i := 0; i < int(tilesX); i++ {
		for j := 0; j < int(tilesY); j++ {
			// Calculate tile position with overlap
			tileX := startX + (float64(i) * tileWidth * 0.95)
			tileY := startY + (float64(j) * tileHeight * 0.95)

			// Create base tile transform
			tileTransform := createTileTransform(tileX, tileY, false)
			transforms = append(transforms, tileTransform)

			// Add 0-1 decorations per tile
			if rand.Float64() < 0.5 { // 50% chance for a decoration
				// Calculate decoration offsets relative to tile
				offsetX := rand.Float64() * (tileWidth * 0.8)  // 80% of tile width
				offsetY := rand.Float64() * (tileHeight * 0.8) // 80% of tile height

				decorTransform := createTileTransform(tileX+offsetX, tileY+offsetY, true)
				transforms = append(transforms, decorTransform)
			}
		}
	}

	// Batch insert transforms
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating transforms: %w", err)
	}

	// Create assets using the created transforms
	for _, transform := range transforms {
		asset := ColonyAssetInsertDTO{
			Colony:    colonyID,
			Transform: transform.ID,
		}

		if transform.ZIndex == 0 {
			asset.AssetCollectionID = BaseGroundTileCollectionID
		} else {
			asset.AssetCollectionID = getRandomDecorationCollectionID()
		}

		assets = append(assets, asset)
	}

	// Batch insert assets
	if err := tx.Create(&assets).Error; err != nil {
		return nil, fmt.Errorf("error creating assets: %w", err)
	}

	// Collect asset IDs
	for _, asset := range assets {
		insertedAssetIDs = append(insertedAssetIDs, int(asset.ID))
	}

	return insertedAssetIDs, nil
}
