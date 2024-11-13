package colony

import (
	"fmt"

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

func InitializeColonyPaths(tx *gorm.DB, colonyID uint32, colonyLocationIDMap map[uint]uint) error {
	paths := []struct {
		LocationA uint
		LocationB uint
	}{
		{uint(getTownHallID()), uint(getHomeID())},
		{uint(getTownHallID()), uint(getCantinaID())},
		{uint(getTownHallID()), uint(getVehicleStorageID())},
		{uint(getTownHallID()), uint(getAgricultureCenterID())},
		{uint(getTownHallID()), uint(getAquiferPlantID())},
		{uint(getTownHallID()), uint(getShieldGeneratorsID())},
		{uint(getHomeID()), uint(getCantinaID())},
		{uint(getVehicleStorageID()), uint(getMiningFacilityID())},
		{uint(getVehicleStorageID()), uint(getRadarDishID())},
		{uint(getAgricultureCenterID()), uint(getAquiferPlantID())},
		{uint(getRadarDishID()), uint(getOuterWallsID())},
		{uint(getRadarDishID()), uint(getSpacePortID())},
	}

	for _, loc := range paths {

		// Use the map to get the ColonyLocation ID
		locationAID, okA := colonyLocationIDMap[loc.LocationA]
		locationBID, okB := colonyLocationIDMap[loc.LocationB]

		if !okA || !okB {
			return fmt.Errorf("Missing location ID: LocationA: %v, LocationB: %v\n", loc.LocationA, loc.LocationB)
		}

		// Insert the path using ColonyLocation IDs
		path := ColonyLocationPath{
			Colony:    colonyID,
			LocationA: uint32(locationAID), // Cast to uint32
			LocationB: uint32(locationBID), // Cast to uint32
		}

		if err := tx.Create(&path).Error; err != nil {
			return err
		}

		// Optionally, insert the reverse path
		reversePath := ColonyLocationPath{
			Colony:    colonyID,
			LocationA: uint32(locationBID), // Cast to uint32
			LocationB: uint32(locationAID), // Cast to uint32
		}

		if err := tx.Create(&reversePath).Error; err != nil {
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
