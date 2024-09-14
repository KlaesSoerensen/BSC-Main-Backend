package api

import "otte_main_backend/src/util"

type LODDetails struct {
	ID             uint32 `json:"id"`
	DetailLevel    uint32 `json:"detailLevel" gorm:"column:detailLevel"`
	GraphicalAsset uint32 `json:"graphicalAsset" gorm:"column:graphicalAsset"`
}

func (l LODDetails) TableName() string {
	return "LOD"
}

type TransformDTO struct {
	ID      uint32  `json:"id"`
	XOffset float32 `json:"xOffset"`
	YOffset float32 `json:"yOffset"`
	ZIndex  uint32  `json:"zIndex"`
	XScale  float32 `json:"xScale"`
	YScale  float32 `json:"yScale"`
}

type MinimizedAssetDTO struct {
	HasLODs bool         `json:"hasLODs"`
	Width   uint32       `json:"width"`
	Height  uint32       `json:"height"`
	LODs    []LODDetails `json:"LODs"`
}

// PlayerDTO represents the data returned for a player's basic information.
type PlayerDTO struct {
	ID                   uint32          `json:"id"`
	ReferenceID          string          `json:"referenceID" gorm:"column:referenceID"`
	IGN                  string          `json:"IGN"`
	Sprite               uint32          `json:"sprite"`
	Achievements         util.PGIntArray `json:"achievements"`
	HasCompletedTutorial bool            `json:"hasCompletedTutorial"`
}

func (p PlayerDTO) TableName() string {
	return "Player"
}
