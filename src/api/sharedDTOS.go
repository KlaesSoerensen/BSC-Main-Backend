package api

type LODDetails struct {
	ID             uint32 `json:"id"`
	DetailLevel    uint32 `json:"detailLevel"`
	GraphicalAsset uint32 `json:"graphicalAsset"`
}

type TransformDTO struct {
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
