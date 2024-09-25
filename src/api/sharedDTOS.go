package api

import (
	"otte_main_backend/src/util"
)

type AchievementModel struct {
	ID          uint32 `json:"id" gorm:"column:id;primaryKey"`
	Description string `json:"description" gorm:"column:description"`
	Icon        uint32 `json:"icon" gorm:"column:icon"`
	Title       string `json:"title" gorm:"column:title"`
}

func (a *AchievementModel) TableName() string {
	return "Achievement"
}

type LODDetails struct {
	ID             uint32 `json:"id"`
	DetailLevel    uint32 `json:"detailLevel" gorm:"column:detailLevel"`
	GraphicalAsset uint32 `json:"graphicalAsset" gorm:"column:graphicalAsset"`
}

func (l LODDetails) TableName() string {
	return "LOD"
}

type LODDetailsDTO struct {
	ID          uint32 `json:"id"`
	DetailLevel uint32 `json:"detailLevel" gorm:"column:detailLevel"`
}

func (l LODDetailsDTO) TableName() string {
	return "LOD"
}

type TransformModel struct {
	ID      uint32  `gorm:"column:id;primaryKey"`
	XScale  float32 `gorm:"column:xScale"`
	YScale  float32 `gorm:"column:yScale"`
	XOffset float32 `gorm:"column:xOffset"`
	YOffset float32 `gorm:"column:yOffset"`
	ZIndex  int     `gorm:"column:zIndex"`
}

func (tm *TransformModel) TableName() string {
	return "Transform"
}

type TransformDTO struct {
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
	FirstName            string          `json:"firstName" gorm:"column:firstName"`
	LastName             string          `json:"lastName" gorm:"column:lastName"`
	Sprite               uint32          `json:"sprite"`
	Achievements         util.PGIntArray `json:"achievements"`
	HasCompletedTutorial bool            `json:"hasCompletedTutorial" gorm:"-"`
}

func (p PlayerDTO) TableName() string {
	return "Player"
}

// PlayerModel represents the database model for a player.
type PlayerModel struct {
	ID           uint32          `json:"id"`
	ReferenceID  string          `json:"referenceID" gorm:"column:referenceID"`
	FirstName    string          `json:"firstName" gorm:"column:firstName"`
	LastName     string          `json:"lastName" gorm:"column:lastName"`
	Sprite       uint32          `json:"sprite"`
	Achievements util.PGIntArray `json:"achievements"`
}

func (p PlayerModel) TableName() string {
	return "Player"
}
