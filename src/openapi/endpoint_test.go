package openapi

import (
	"otte_main_backend/src/meta"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestApiDefinition_Add(t *testing.T) {
	var cases = []struct {
		name    string
		url     string
		handler openapi.HandlerFunction
		options *openapi.EndpointOptions
	}{
		{
			name: "Does it fill default values correctly with empty options?",
			url:  "Test1",
			handler: func(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
				return c.SendString("Hello, World!")
			},
			options: &EndpointOptions{},
		},
		{
			name: "Does it fill default values correctly with no options?",
			url:  "Test1",
			handler: func(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
				return c.SendString("Hello, World!")
			},
			options: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			apiDef := New(nil)
			apiDef.Add(c.url, c.handler, c.options)
			if apiDef.endpoints[0].UrlExtension != c.url {
				t.Errorf("UrlExtension is not preserved correctly")
			}
			if apiDef.endpoints[0].AuthRequired != AuthDefault {
				t.Errorf("AuthRequired is not defaulted to the right type")
			}
			if apiDef.endpoints[0].Method != HTTPGet {
				t.Errorf("Method is not defaulted to the right type")
			}
			if apiDef.endpoints[0].Handler == nil {
				t.Errorf("Handler is not preserved correctly")
			}
		})
	}

}
