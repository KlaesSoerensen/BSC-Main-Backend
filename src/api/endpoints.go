package api

import (
	"otte_main_backend/src/openapi"
)

func ApplyEndpoints(apiDef *openapi.ApiDefinition) error {
	applyHealthApi(apiDef)

	return nil
}
