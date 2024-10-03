package colony

import (
	"gorm.io/gorm"
)

type ColonyLocationPath struct {
	ID        uint32 `gorm:"column:id;primaryKey"`
	Colony    uint32 `gorm:"column:colony"`
	LocationA uint32 `gorm:"column:locationA"`
	LocationB uint32 `gorm:"column:locationB"`
}

func (ColonyLocationPath) TableName() string {
	return "ColonyLocationPath"
}

func InitializeColonyPaths(tx *gorm.DB, colonyID uint32) error {
	paths := []ColonyLocationPath{
		// Central Hub: Town Hall
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getHomeID()},
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getCantinaID()},
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getVehicleStorageID()},
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getAgricultureCenterID()},
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getAquiferPlantID()},
		{Colony: colonyID, LocationA: getTownHallID(), LocationB: getShieldGeneratorsID()},
		// Residential Area
		{Colony: colonyID, LocationA: getHomeID(), LocationB: getCantinaID()},
		// Industrial and Defensive Area
		{Colony: colonyID, LocationA: getVehicleStorageID(), LocationB: getMiningFacilityID()},
		{Colony: colonyID, LocationA: getVehicleStorageID(), LocationB: getRadarDishID()},
		// Agricultural Zone
		{Colony: colonyID, LocationA: getAgricultureCenterID(), LocationB: getAquiferPlantID()},
		// Defensive Perimeter and Endpoints
		{Colony: colonyID, LocationA: getRadarDishID(), LocationB: getOuterWallsID()},
		{Colony: colonyID, LocationA: getRadarDishID(), LocationB: getSpacePortID()},
	}

	// Create duplicate entries for each path to represent both directions
	duplicatedPaths := make([]ColonyLocationPath, 0, len(paths)*2)
	for _, path := range paths {
		// Original direction
		duplicatedPaths = append(duplicatedPaths, path)
		// Reverse direction
		duplicatedPaths = append(duplicatedPaths, ColonyLocationPath{
			Colony:    path.Colony,
			LocationA: path.LocationB,
			LocationB: path.LocationA,
		})
	}

	// Insert the paths into the database
	for _, path := range duplicatedPaths {
		if err := tx.Create(&path).Error; err != nil {
			return err
		}
	}

	return nil
}

func getTownHallID() uint32          { return 40 }
func getHomeID() uint32              { return 30 }
func getCantinaID() uint32           { return 90 }
func getVehicleStorageID() uint32    { return 80 }
func getAgricultureCenterID() uint32 { return 70 }
func getAquiferPlantID() uint32      { return 60 }
func getMiningFacilityID() uint32    { return 110 }
func getOuterWallsID() uint32        { return 10 }
func getShieldGeneratorsID() uint32  { return 50 }
func getRadarDishID() uint32         { return 100 }
func getSpacePortID() uint32         { return 20 }
