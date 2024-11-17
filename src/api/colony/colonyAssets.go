// Package colony provides functionality for managing and generating graphical assets
// in a colony-based game or simulation environment. It handles the creation and
// placement of various visual elements including ground tiles, walls, glass panels,
// and decorative elements.
package colony

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// GraphicalAsset represents a visual element that can be placed in the colony.
// It contains metadata about the asset's dimensions and purpose.
type GraphicalAsset struct {
	ID      uint32 `gorm:"column:id;primaryKey"` // Unique identifier for the asset
	Alias   string `gorm:"column:alias"`         // Human-readable name for the asset
	Type    string `gorm:"column:type"`          // Category or type of the asset
	UseCase string `gorm:"column:useCase"`       // Intended usage context
	Width   int    `gorm:"column:width"`         // Width in pixels
	Height  int    `gorm:"column:height"`        // Height in pixels
}

// ColonyAssetInsertDTO represents the data structure used for inserting
// new assets into a colony, including their positioning and relationship
// to asset collections.
type ColonyAssetInsertDTO struct {
	ID                uint32 `gorm:"column:id;primaryKey"`                          // Unique identifier
	AssetCollectionID uint32 `json:"assetCollection" gorm:"column:assetCollection"` // Reference to collection
	Transform         uint   `json:"transform" gorm:"column:transform"`             // Position and scale data
	Colony            uint32 `json:"colony" gorm:"column:colony"`                   // Associated colony ID
}

// TableName returns the database table name for GraphicalAsset
func (GraphicalAsset) TableName() string {
	return "GraphicalAsset"
}

// TableName returns the database table name for ColonyAssetInsertDTO
func (cai *ColonyAssetInsertDTO) TableName() string {
	return "ColonyAsset"
}

// createTileTransform generates a new Transform for a tile with specified position
// and layer information.
// Parameters:
//   - x: X-coordinate position
//   - y: Y-coordinate position
//   - isDecoration: determines if the tile should be placed on the decoration layer
//
// Returns:
//   - Transform: A new transform with the specified position and scaling
func createTileTransform(x, y float64, isDecoration bool) Transform {
	// Set the z-index based on whether this is a decoration
	zIndex := 0
	if isDecoration {
		zIndex = 1
	}

	// Create and return a new transform with default scaling
	return Transform{
		XScale:  1,
		YScale:  1,
		XOffset: x,
		YOffset: y,
		ZIndex:  zIndex,
	}
}

// createRandomOffset generates random X and Y offsets within a cell
// for decoration placement.
// Parameters:
//   - cellSize: The size of the cell to generate offsets within
//
// Returns:
//   - float64: X offset
//   - float64: Y offset
func createRandomOffset(cellSize float64) (float64, float64) {
	// Calculate maximum allowed offset as 80% of cell size
	maxOffset := cellSize * 0.8

	// Generate and return random offsets between -maxOffset/2 and +maxOffset/2
	return (rand.Float64() - 0.5) * maxOffset, (rand.Float64() - 0.5) * maxOffset
}

// createRandomScale generates a random scale value between specified bounds.
// Parameters:
//   - minScale: Minimum scale value
//   - maxScale: Maximum scale value
//
// Returns:
//   - float64: Generated scale value between minScale and maxScale
func createRandomScale(minScale, maxScale float64) float64 {
	// Generate random scale within the specified range
	return minScale + rand.Float64()*(maxScale-minScale)
}

// createDecorationTransform generates a Transform for decorative elements
// with custom scaling and positioning.
// Parameters:
//   - x: X-coordinate position
//   - y: Y-coordinate position
//   - scale: Scale factor for the decoration
//   - isDecoration: determines if this should be placed on the decoration layer
//
// Returns:
//   - Transform: A new transform with the specified position, scale, and z-index
func createDecorationTransform(x, y float64, scale float64, isDecoration bool) Transform {
	// Set the z-index based on whether this is a decoration
	zIndex := 0
	if isDecoration {
		zIndex = 1
	}

	// Create and return a new transform with custom scaling
	return Transform{
		XScale:  scale,
		YScale:  scale,
		XOffset: x,
		YOffset: y,
		ZIndex:  zIndex,
	}
}

// createTiles generates the ground tiles for the colony within the specified bounds.
// Parameters:
//   - tx: Database transaction
//   - colonyID: ID of the colony
//   - baseTile: Reference tile asset
//   - expandedBoundingBox: Area to fill with tiles
//   - globalYOffsetWall: Global Y-axis offset for proper layering
//
// Returns:
//   - []ColonyAssetInsertDTO: Slice of created assets
//   - error: Any error encountered during creation
func createTiles(tx *gorm.DB, colonyID uint32, baseTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	// Calculate adjusted tile dimensions with 90% of original size to create slight overlap
	adjustedTileWidth := float64(baseTile.Width) * 0.9
	adjustedTileHeight := float64(baseTile.Height) * 0.9

	// Calculate the total area to be covered
	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	deltaY := expandedBoundingBox.MaxY - expandedBoundingBox.MinY

	// Estimate the number of tiles needed
	estimatedSize := int(math.Floor((deltaX / adjustedTileWidth) * (deltaY / adjustedTileHeight)))

	// Initialize slice to store transforms
	transforms := make([]Transform, 0, estimatedSize)

	// Generate transforms for each tile position
	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		for y := expandedBoundingBox.MinY; y < expandedBoundingBox.MaxY; y += adjustedTileHeight {
			tileTransform := createTileTransform(x, y+globalYOffsetWall, false)
			transforms = append(transforms, tileTransform)
		}
	}

	// Save all transforms to the database
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating transforms: %s", err.Error())
	}

	// Create colony assets using the generated transforms
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

// createWallTiles generates wall tiles along the colony boundary.
// Parameters:
//   - tx: Database transaction
//   - colonyID: ID of the colony
//   - wallTile: Wall tile asset
//   - expandedBoundingBox: Area to place walls within
//   - globalYOffsetWall: Global Y-axis offset for proper layering
//
// Returns:
//   - []ColonyAssetInsertDTO: Slice of created wall assets
//   - error: Any error encountered during creation
func createWallTiles(tx *gorm.DB, colonyID uint32, wallTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	// Calculate adjusted tile dimensions with 90% of original size
	adjustedTileWidth := float64(wallTile.Width) * 0.9
	adjustedTileHeight := float64(wallTile.Height) * 0.9

	// Calculate the total width to be covered
	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX

	// Estimate the number of wall tiles needed
	estimatedSize := int(math.Floor(deltaX / adjustedTileWidth))

	// Initialize slice to store transforms
	transforms := make([]Transform, 0, estimatedSize)

	// Calculate wall Y position relative to ground position
	wallYPosition := (globalYOffsetWall / 2) - adjustedTileHeight - 225

	// Generate transforms for each wall tile position
	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		wallTransform := createTileTransform(x, wallYPosition, false)
		transforms = append(transforms, wallTransform)
	}

	// Save all transforms to the database
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating wall transforms: %s", err.Error())
	}

	// Create colony assets using the generated transforms
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

// createRandomDecorations generates randomly placed decorative elements within the colony.
// Parameters:
//   - tx: Database transaction
//   - colonyID: ID of the colony
//   - baseTile: Reference tile for sizing
//   - expandedBoundingBox: Area to place decorations within
//   - globalYOffsetWall: Global Y-axis offset for proper layering
//
// Returns:
//   - []ColonyAssetInsertDTO: Slice of created decoration assets
//   - error: Any error encountered during creation
func createRandomDecorations(tx *gorm.DB, colonyID uint32, baseTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	// Initialize random seed (deprecated - issue?)
	rand.Seed(time.Now().UnixNano())

	// Calculate cell dimensions for decoration placement
	adjustedTileWidth := float64(baseTile.Width) * 0.9
	adjustedTileHeight := float64(baseTile.Height) * 0.9
	cellWidth := adjustedTileWidth / 3
	cellHeight := adjustedTileHeight / 3

	// Calculate total area and maximum possible decorations
	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	deltaY := expandedBoundingBox.MaxY - expandedBoundingBox.MinY
	maxDecorations := int(math.Floor((deltaX / cellWidth) * (deltaY / cellHeight)))

	// Estimate actual number of decorations (1/3 of max)
	estimatedSize := maxDecorations / 3

	// Initialize slices for transforms and assets
	transforms := make([]Transform, 0, estimatedSize)
	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)

	// Generate random decorations
	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += cellWidth {
		for y := expandedBoundingBox.MinY; y < expandedBoundingBox.MaxY; y += cellHeight {
			// 33% chance to place a decoration
			if rand.Float64() <= 0.33 {
				// Generate random position offsets and scale
				offsetX, offsetY := createRandomOffset(cellWidth)
				finalX := x + offsetX
				finalY := y + offsetY + globalYOffsetWall
				scale := createRandomScale(0.5, 1.0)

				// Create transform for the decoration
				decorTransform := createDecorationTransform(finalX, finalY, scale, true)
				transforms = append(transforms, decorTransform)
			}
		}
	}

	// Save all transforms to the database
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating decoration transforms: %s", err.Error())
	}

	// Create colony assets with random decoration types
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

// createGlassTiles generates glass tiles along the top of the colony walls.
// Parameters:
//   - tx: Database transaction
//   - colonyID: ID of the colony
//   - glassTile: Glass tile asset
//   - expandedBoundingBox: Area to place glass within
//   - globalYOffsetWall: Global Y-axis offset for proper layering
//
// Returns:
//   - []ColonyAssetInsertDTO: Slice of created glass assets
//   - error: Any error encountered during creation
func createGlassTiles(tx *gorm.DB, colonyID uint32, glassTile *GraphicalAsset, expandedBoundingBox *BoundingBox, globalYOffsetWall float64) ([]ColonyAssetInsertDTO, error) {
	// Calculate adjusted tile dimensions with 99% of original size
	adjustedTileWidth := float64(glassTile.Width) * 0.99
	adjustedTileHeight := float64(glassTile.Height) * 0.99

	// Calculate total width to be covered
	deltaX := expandedBoundingBox.MaxX - expandedBoundingBox.MinX
	// Double the estimated size for two rows
	estimatedSize := int(math.Floor(deltaX/adjustedTileWidth)) * 2

	// Initialize slice for transforms
	transforms := make([]Transform, 0, estimatedSize)

	// Calculate base glass tile Y position relative to wall position
	wallYPosition := (globalYOffsetWall / 2) - adjustedTileHeight - 125

	// Y positions for both rows of glass
	// First row (lower)
	glassYPosition1 := wallYPosition - adjustedTileHeight + 100 + 100
	// Second row (upper) - offset by adjusted tile height with a small gap
	glassYPosition2 := glassYPosition1 - adjustedTileHeight

	// Generate transforms for each glass tile position - first row
	for x := expandedBoundingBox.MinX; x < expandedBoundingBox.MaxX; x += adjustedTileWidth {
		// Create transform for lower row
		glassTransform1 := createTileTransform(x, glassYPosition1, false)
		transforms = append(transforms, glassTransform1)

		// Create transform for upper row
		glassTransform2 := createTileTransform(x, glassYPosition2, false)
		transforms = append(transforms, glassTransform2)
	}

	// Save all transforms to the database
	if err := tx.Create(&transforms).Error; err != nil {
		return nil, fmt.Errorf("error creating glass transforms: %s", err.Error())
	}

	// Create colony assets using the generated transforms
	assets := make([]ColonyAssetInsertDTO, 0, estimatedSize)
	for _, transform := range transforms {
		assets = append(assets, ColonyAssetInsertDTO{
			Colony:            colonyID,
			Transform:         transform.ID,
			AssetCollectionID: 10035,
		})
	}

	return assets, nil
}

// InsertColonyAssets is the main entry point for creating all visual elements of a colony.
// It coordinates the creation of ground tiles, walls, glass panels, and decorations,
// ensuring proper layering and positioning of all elements.
//
// The assets are created in the following order to ensure proper layering:
// 1. Ground tiles (bottom layer)
// 2. Glass tiles (behind walls)
// 3. Wall tiles (in front of glass)
// 4. Decorations (top layer)
//
// Parameters:
//   - tx: Database transaction
//   - colonyID: ID of the colony to create assets for
//   - boundingBox: Defines the area of the colony
//
// Returns:
//   - []int: Slice of inserted asset IDs
//   - error: Any error encountered during the creation process
func InsertColonyAssets(tx *gorm.DB, colonyID uint32, boundingBox *BoundingBox) ([]int, error) {
	// Fetch the base tile asset from the database
	var baseTile GraphicalAsset
	if err := tx.Where("id = ?", 8001).First(&baseTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching base tile asset: %w", err)
	}

	// Fetch the wall tile asset from the database
	var wallTile GraphicalAsset
	if err := tx.Where("id = ?", 8034).First(&wallTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching wall tile asset: %w", err)
	}

	// Fetch the glass tile asset from the database
	var glassTile GraphicalAsset
	if err := tx.Where("id = ?", 8035).First(&glassTile).Error; err != nil {
		return nil, fmt.Errorf("error fetching glass tile asset: %w", err)
	}

	// Calculate expanded bounding box to ensure coverage beyond visible area
	expandedBoundingBox := BoundingBox{
		MinX: boundingBox.MinX - (1920/2 + 200), // Expand by half screen width
		MaxX: boundingBox.MaxX + (1920/2 + 200),
		MinY: boundingBox.MinY - (1080 / 2), // Expand by half screen height
		MaxY: boundingBox.MaxY + (1080 / 2),
	}

	// Calculate global Y offset for wall positioning
	globalYOffsetWall := (boundingBox.MinY - expandedBoundingBox.MinY) * 2

	// Create ground tiles (bottom layer)
	tileAssets, err := createTiles(tx, colonyID, &baseTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating tiles: %w", err)
	}

	// Create glass tiles (middle layer, behind walls)
	glassAssets, err := createGlassTiles(tx, colonyID, &glassTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating glass tiles: %w", err)
	}

	// Create wall tiles (middle layer, in front of glass)
	wallAssets, err := createWallTiles(tx, colonyID, &wallTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating wall tiles: %w", err)
	}

	// Create decorative elements (top layer)
	decorAssets, err := createRandomDecorations(tx, colonyID, &baseTile, &expandedBoundingBox, globalYOffsetWall)
	if err != nil {
		return nil, fmt.Errorf("error creating decorations: %w", err)
	}

	// Combine all assets in the correct layering order
	allAssets := append(tileAssets, glassAssets...) // Ground tiles and glass tiles first
	allAssets = append(allAssets, wallAssets...)    // Wall tiles next
	allAssets = append(allAssets, decorAssets...)   // Decorations on top

	// Save all assets to the database
	if err := tx.Create(&allAssets).Error; err != nil {
		return nil, fmt.Errorf("error creating assets: %w", err)
	}

	// Create slice of inserted asset IDs for return
	insertedAssetIDs := make([]int, len(allAssets))
	for i, asset := range allAssets {
		insertedAssetIDs[i] = int(asset.ID)
	}

	return insertedAssetIDs, nil
}
