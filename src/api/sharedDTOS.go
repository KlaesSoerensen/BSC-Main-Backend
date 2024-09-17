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
	XOffset float32 `json:"xOffset" gorm:"column:xOffset"`
	YOffset float32 `json:"yOffset" gorm:"column:yOffset"`
	ZIndex  uint32  `json:"zIndex" gorm:"column:zIndex"`
	XScale  float32 `json:"xScale" gorm:"column:xScale"`
	YScale  float32 `json:"yScale" gorm:"column:yScale"`
}

func (t *TransformDTO) TableName() string {
	return "Transform"
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
	IGN                  string          `json:"IGN" gorm:"column:IGN"`
	Sprite               uint32          `json:"sprite"`
	Achievements         util.PGIntArray `json:"achievements"`
	HasCompletedTutorial bool            `json:"hasCompletedTutorial" gorm:"-"`
}

func (p PlayerDTO) TableName() string {
	return "Player"
}
