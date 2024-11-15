package colony

import (
	"otte_main_backend/src/meta"

	"gorm.io/gorm"
)

type Transform struct {
	ID      uint    `gorm:"column:id;primaryKey"`
	XScale  float64 `gorm:"column:xScale"`
	YScale  float64 `gorm:"column:yScale"`
	XOffset float64 `gorm:"column:xOffset"`
	YOffset float64 `gorm:"column:yOffset"`
	ZIndex  int     `gorm:"column:zIndex"`
}

func (Transform) TableName() string {
	return "Transform"
}

func createTransform(xScale float64, yScale float64, xOffset float64, yOffset float64, zIndex int) Transform {
	topLevelScaleScalar := .75
	topLevelDistanceScalar := float64(2)
	return Transform{
		XScale:  xScale * topLevelScaleScalar,
		YScale:  yScale * topLevelScaleScalar,
		XOffset: xOffset * topLevelDistanceScalar,
		YOffset: yOffset * topLevelDistanceScalar,
		ZIndex:  zIndex,
	}
}

type BoundingBox struct {
	MinX float64
	MaxX float64
	MinY float64
	MaxY float64
}

func InsertLocationTransforms(appContext *meta.ApplicationContext, tx *gorm.DB) (map[string]uint, error, *BoundingBox) {
	// Positions are scaled to act as if 2048 x 1080
	transforms := []Transform{
		createTransform(1, 1, 650, 400, 100),  // Town Hall
		createTransform(1, 1, 400, 280, 100),  // Cantina
		createTransform(1, 1, 220, 380, 100),  // Home
		createTransform(1, 1, 850, 220, 100),  // Aquifer Plant
		createTransform(1, 1, 600, 580, 100),  // Shield Generator
		createTransform(1, 1, 1020, 420, 100), // Vehicle Storage
		createTransform(1, 1, 1450, 300, 100), // Radar Dish
		createTransform(1, 1, 1400, 500, 100), // Mining Facility
		createTransform(1, 1, 1750, 280, 100), // Outer Walls
		createTransform(1, 1, 1800, 450, 100), // Space Port
		createTransform(1, 1, 620, 150, 100),  // Agriculture Center
	}

	boundingBox := BoundingBox{}
	for _, transform := range transforms {
		if transform.XOffset < boundingBox.MinX {
			boundingBox.MinX = transform.XOffset
		}
		if transform.XOffset > boundingBox.MaxX {
			boundingBox.MaxX = transform.XOffset
		}
		if transform.YOffset < boundingBox.MinY {
			boundingBox.MinY = transform.YOffset
		}
		if transform.YOffset > boundingBox.MaxY {
			boundingBox.MaxY = transform.YOffset
		}
	}

	transformIDs := make(map[string]uint)
	locationNames := []string{
		"LOCATION.TOWN_HALL.NAME",
		"LOCATION.CANTINA.NAME",
		"LOCATION.HOME.NAME",
		"LOCATION.AQUIFER_PLANT.NAME",
		"LOCATION.SHIELD_GENERATOR.NAME",
		"LOCATION.VEHICLE_STORAGE.NAME",
		"LOCATION.RADAR_DISH.NAME",
		"LOCATION.MINING_FACILITY.NAME",
		"LOCATION.OUTER_WALLS.NAME",
		"LOCATION.SPACE_PORT.NAME",
		"LOCATION.AGRICULTURE_CENTER.NAME",
	}

	for i, transform := range transforms {
		if err := tx.Create(&transform).Error; err != nil {
			return nil, err, nil
		}
		transformIDs[locationNames[i]] = transform.ID
	}

	return transformIDs, nil, &boundingBox
}
