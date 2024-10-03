package colony

import (
	"otte_main_backend/src/meta"

	"gorm.io/gorm"
)

type ColonyLocation struct {
	ID        uint `gorm:"column:id;primaryKey"`
	Colony    uint `gorm:"column:colony"`
	Location  uint `gorm:"column:location"`
	Transform uint `gorm:"column:transform"`
	Level     int  `gorm:"column:level"`
}

// TableName specifies the table name for this struct
func (ColonyLocation) TableName() string {
	return "ColonyLocation"
}

func InsertColonyLocations(appContext *meta.ApplicationContext, tx *gorm.DB, colonyID uint, transformIDs map[string]uint) error {
	locations := []struct {
		Name  string
		ID    uint
		Level int
	}{
		{"LOCATION.TOWN_HALL.NAME", 40, 1},
		{"LOCATION.CANTINA.NAME", 90, 1},
		{"LOCATION.HOME.NAME", 30, 1},
		{"LOCATION.AQUIFER_PLANT.NAME", 60, 1},
		{"LOCATION.SHIELD_GENERATOR.NAME", 50, 1},
		{"LOCATION.VEHICLE_STORAGE.NAME", 80, 1},
		{"LOCATION.RADAR_DISH.NAME", 100, 1},
		{"LOCATION.MINING_FACILITY.NAME", 110, 1},
		{"LOCATION.OUTER_WALLS.NAME", 10, 1},
		{"LOCATION.SPACE_PORT.NAME", 20, 1},
	}

	for _, loc := range locations {
		colonyLocation := ColonyLocation{
			Colony:    colonyID,
			Location:  loc.ID,
			Transform: transformIDs[loc.Name],
			Level:     loc.Level,
		}
		if err := tx.Create(&colonyLocation).Error; err != nil {
			return err
		}
	}

	return nil
}
