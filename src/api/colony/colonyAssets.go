package colony

import (
	"fmt"
	"math"
	"math/rand"
	"time"

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

func createRandomOffset(cellSize float64) (float64, float64) {
	maxOffset := cellSize * 0.8
	return (rand.Float64() - 0.5) * maxOffset, (rand.Float64() - 0.5) * maxOffset
}

func createRandomScale(minScale, maxScale float64) float64 {
	return minScale + rand.Float64()*(maxScale-minScale)
}

func createDecorationTransform(x, y float64, scale float64, isDecoration bool) Transform {
	zIndex := 0
	if isDecoration {
		zIndex = 1
	}
	return Transform{
		XScale:  scale,
		YScale:  scale,
		XOffset: x,
		YOffset: y,
		ZIndex:  zIndex,
	}
}

func createTiles(tx *gorm.DB, colonyID uint32, baseTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	adjustedTileWidth := float64(baseTile.Width) * 0.9
	adjustedTileHeight := float64(baseTile.Height) * 0.9

	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	deltaY := expandedBoundingBox.MaxY - expandedBoundingBox.MinY
	estimatedSize := int(math.Floor((deltaX / adjustedTileWidth) * (deltaY / adjustedTileHeight)))

	transforms := make([]Transform, 0, estimatedSize)

	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		for y := expandedBoundingBox.MinY; y < expandedBoundingBox.MaxY; y += adjustedTileHeight {
			tileTransform := createTileTransform(x, y+globalYOffsetWall, false)
			transforms = append(transforms, tileTransform)
		}
	}

	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating transforms: %s", err.Error())
	}

	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)
	for _, transform := range transforms {
		assets = append(assets, ColonyAssetInsertDTO{
			Colony:            colonyID,
			Transform:         transform.ID,
			AssetCollectionID: 10001,
		})
	}

	return assets, nil
}

func createWallTiles(tx *gorm.DB, colonyID uint32, wallTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	adjustedTileWidth := float64(wallTile.Width) * 0.9
	adjustedTileHeight := float64(wallTile.Height) * 0.9

	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	estimatedSize := int(math.Floor(deltaX / adjustedTileWidth))

	transforms := make([]Transform, 0, estimatedSize)

	// The wall should be positioned where the ground starts (globalYOffsetWall)
	// but considering that globalYOffsetWall is doubled, we use it without doubling
	// and subtract the wall height to place it above
	wallYPosition := (globalYOffsetWall / 2) - adjustedTileHeight - 125

	fmt.Printf("Wall attributes - Width: %d, Height: %d\n", wallTile.Width, wallTile.Height)
	fmt.Printf("Adjusted wall height: %f\n", adjustedTileHeight)
	fmt.Printf("globalYOffsetWall: %f\n", globalYOffsetWall)
	fmt.Printf("New wallYPosition: %f\n", wallYPosition)

	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		wallTransform := createTileTransform(x, wallYPosition, false)
		transforms = append(transforms, wallTransform)
	}

	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating wall transforms: %s", err.Error())
	}

	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)
	for _, transform := range transforms {
		assets = append(assets, ColonyAssetInsertDTO{
			Colony:            colonyID,
			Transform:         transform.ID,
			AssetCollectionID: 10034,
		})
	}

	return assets, nil
}

func createRandomDecorations(tx *gorm.DB, colonyID uint32, baseTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	rand.Seed(time.Now().UnixNano())

	adjustedTileWidth := float64(baseTile.Width) * 0.9
	adjustedTileHeight := float64(baseTile.Height) * 0.9
	cellWidth := adjustedTileWidth / 3
	cellHeight := adjustedTileHeight / 3

	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	deltaY := expandedBoundingBox.MaxY - expandedBoundingBox.MinY
	maxDecorations := int(math.Floor((deltaX / cellWidth) * (deltaY / cellHeight)))

	estimatedSize := maxDecorations / 3

	transforms := make([]Transform, 0, estimatedSize)
	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)

	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += cellWidth {
		for y := expandedBoundingBox.MinY; y < expandedBoundingBox.MaxY; y += cellHeight {
			if rand.Float64() <= 0.33 {
				offsetX, offsetY := createRandomOffset(cellWidth)
				finalX := x + offsetX
				finalY := y + offsetY + globalYOffsetWall
				scale := createRandomScale(0.5, 1.0)
				decorTransform := createDecorationTransform(finalX, finalY, scale, true)
				transforms = append(transforms, decorTransform)
			}
		}
	}

	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating decoration transforms: %s", err.Error())
	}

	for _, transform := range transforms {
		collectionID := uint32(10002 + rand.Intn(31))
		assets = append(assets, ColonyAssetInsertDTO{
			Colony:            colonyID,
			Transform:         transform.ID,
			AssetCollectionID: collectionID,
		})
	}

	return assets, nil
}

func InsertColonyAssets(tx *gorm.DB, colonyID uint32, boundingBox *BoundingBox) ([]int, error) {
	var baseTile GraphicalAsset
	if err := tx.Where("id = ?", 8001).First(&baseTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching base tile asset: %w", err)
	}

	var wallTile GraphicalAsset
	if err := tx.Where("id = ?", 8034).First(&wallTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching wall tile asset: %w", err)
	}

	expandedBoundingBox := BoundingBox{
		MinX: boundingBox.MinX - (1920 / 2),
		MaxX: boundingBox.MaxX + (1920 / 2),
		MinY: boundingBox.MinY - (1080 / 2),
		MaxY: boundingBox.MaxY + (1080 / 2),
	}

	globalYOffsetWall := (boundingBox.MinY - expandedBoundingBox.MinY) * 2

	tileAssets, err := createTiles(tx, colonyID, &baseTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating tiles: %w", err)
	}

	wallAssets, err := createWallTiles(tx, colonyID, &wallTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating wall tiles: %w", err)
	}

	decorAssets, err := createRandomDecorations(tx, colonyID, &baseTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating decorations: %w", err)
	}

	allAssets := append(tileAssets, wallAssets...)
	allAssets = append(allAssets, decorAssets...)

	if err := tx.Create(&allAssets).Error; err != nil {
		return nil, fmt.Errorf("error creating assets: %w", err)
	}

	insertedAssetIDs := make([]int, len(allAssets))
	for i, asset := range allAssets {
		insertedAssetIDs[i] = int(asset.ID)
	}

	return insertedAssetIDs, nil
}
