package colony

import (
	"fmt"
	"math"

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

func createTiles(tx *gorm.DB, colonyID uint32, baseTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]Transform, []ColonyAssetInsertDTO, error) {
	// The current base tile is 663x663, subtracting for inverse padding (feathered edges by 5% on each side)
	// should result in placements of 596.7x596.7
	adjustedTileWidth := float64(baseTile.Width) * 0.9
	adjustedTileHeight := float64(baseTile.Height) * 0.9

	// Pre-allocate slices with approximate size
	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	deltaY := expandedBoundingBox.MaxY - expandedBoundingBox.MinY
	estimatedSize := int(math.Floor((deltaX / adjustedTileWidth) * (deltaY / adjustedTileHeight)))

	transforms := make([]Transform, 0, estimatedSize)

	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		for y := expandedBoundingBox.MinY; y < expandedBoundingBox.MaxY; y += adjustedTileHeight {
			// Create base tile transform
			tileTransform := createTileTransform(x, y+globalYOffsetWall, false)
			transforms = append(transforms, tileTransform)
		}
	}

	// Batch insert transforms
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, nil, fmt.Errorf("error creating transforms: %s", err.Error())
	}

	// Create assets using the created transforms
	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)
	for _, transform := range transforms {
		assets = append(assets, ColonyAssetInsertDTO{
			Colony:            colonyID,
			Transform:         transform.ID,
			AssetCollectionID: 10001,
		})
	}

	return transforms, assets, nil
}

func InsertColonyAssets(tx *gorm.DB, colonyID uint32, boundingBox *BoundingBox) ([]int, error) {
	// Get the base tile asset metadata
	var baseTile GraphicalAsset
	if err := tx.Where("id = ?", 8001).First(&baseTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching base tile asset: %w", err)
	}

	// Everything is normalized to expecting a 2K screen size (1920x1080)
	// so expanding by half that in all directions should assure that when standing
	// at the outermost location, the "edge" is not seen.
	expandedBoundingBox := BoundingBox{
		MinX: boundingBox.MinX - (1920 / 2),
		MaxX: boundingBox.MaxX + (1920 / 2),
		MinY: boundingBox.MinY - (1080 / 2),
		MaxY: boundingBox.MaxY + (1080 / 2),
	}

	globalYOffsetWall := (boundingBox.MinY - expandedBoundingBox.MinY) * 2

	_, assets, err := createTiles(tx, colonyID, &baseTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating tiles: %w", err)
	}

	// Batch insert assets
	if err := tx.Create(&assets).Error; err != nil {
		return nil, fmt.Errorf("error creating assets: %w", err)
	}

	// Collect asset IDs
	insertedAssetIDs := make([]int, len(assets))
	for i, asset := range assets {
		insertedAssetIDs[i] = int(asset.ID)
	}

	return insertedAssetIDs, nil
}
