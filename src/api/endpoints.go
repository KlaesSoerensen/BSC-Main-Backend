package api

import (
	"otte_main_backend/src/openapi"
)

type ServiceConstants struct {
	ServicePort int64
	DDH         string
	AuthToken   string
}

func ApplyEndpoints(apiDef *openapi.ApiDefinition) error {
	applyHealthApi(apiDef)

	return nil
}
