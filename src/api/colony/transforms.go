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

func InsertTransforms(appContext *meta.ApplicationContext, tx *gorm.DB) (map[string]uint, error) {
	topLevelLocationScalar := .75
	transforms := []Transform{
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 400, YOffset: 400, ZIndex: 100},  // Town Hall
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 250, YOffset: 300, ZIndex: 100},  // Cantina
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 100, YOffset: 400, ZIndex: 100},  // Home
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 400, YOffset: 200, ZIndex: 100},  // Aquifer Plant
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 400, YOffset: 600, ZIndex: 100},  // Shield Generator
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 700, YOffset: 400, ZIndex: 100},  // Vehicle Storage
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 1000, YOffset: 300, ZIndex: 100}, // Radar Dish
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 1000, YOffset: 500, ZIndex: 100}, // Mining Facility
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 1300, YOffset: 300, ZIndex: 100}, // Outer Walls
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 1300, YOffset: 500, ZIndex: 100}, // Space Port
		{XScale: 1 * topLevelLocationScalar, YScale: 1 * topLevelLocationScalar, XOffset: 400, YOffset: 0, ZIndex: 100},    // Agriculture Center
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
			return nil, err
		}
		transformIDs[locationNames[i]] = transform.ID
	}

	return transformIDs, nil
}
