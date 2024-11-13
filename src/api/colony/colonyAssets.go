package colony

import "gorm.io/gorm"

type ColonyAssetInsertDTO struct {
	AssetCollectionID uint32    `json:"assetCollection" gorm:"column:assetCollection"`
	Transform         Transform `json:"transform" gorm:"column:transform"`
}

func (cai *ColonyAssetInsertDTO) TableName() string {
	return "ColonyAsset"
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

func InsertColonyAssets(tx *gorm.DB, colonyID uint32, boundingBox *BoundingBox) error {
	//2048 x 1080, padding: 1024 x 540
	// 10001, 400 x 400
	colonyAssets := []*ColonyAssetInsertDTO{
		GenerateColonyAsset(10001, 650, 400, colonyID),
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