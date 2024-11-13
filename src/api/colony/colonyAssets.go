package colony

import "gorm.io/gorm"

type ColonyAssetInsertDTO struct {
	AssetCollectionID uint32    `json:"assetCollection" gorm:"column:assetCollection"`
	Transform         Transform `json:"transform" gorm:"column:transform"`
}

func GenerateColonyAsset(assetCollectionID uint32, x float64, y float64, colonyID uint32) *ColonyAssetInsertDTO {
	return &ColonyAssetInsertDTO{
		AssetCollectionID: assetCollectionID,
		Transform: Transform{
			XOffset: x,
			YOffset: y,
			ZIndex:  1,
			XScale:  1,
			YScale:  1,
		},
	}
}

func InsertColonyAssets(tx *gorm.DB, colonyID uint32) error {
	colonyAssets := []*ColonyAssetInsertDTO{
		GenerateColonyAsset(1, 650, 400, colonyID),
		GenerateColonyAsset(1, 250, 400, colonyID),
		GenerateColonyAsset(1, -150, 400, colonyID),
	}

	for _, asset := range colonyAssets {
		if err := tx.Create(asset).Error; err != nil {
			return err
		}
	}
	return nil
}
