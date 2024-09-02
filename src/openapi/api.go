package openapi

import (
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/config"
	"otte_main_backend/src/meta"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type HTTPMethod = string
type AuthMethod = string

const (
	AuthNone    AuthMethod = "NONE"
	AuthDefault AuthMethod = "DEFAULT"
)

const (
	HTTPGet     HTTPMethod = "GET"     // RFC 7231, 4.3.1
	HTTPHead    HTTPMethod = "HEAD"    // RFC 7231, 4.3.2
	HTTPPost    HTTPMethod = "POST"    // RFC 7231, 4.3.3
	HTTPPut     HTTPMethod = "PUT"     // RFC 7231, 4.3.4
	HTTPPatch   HTTPMethod = "PATCH"   // RFC 5789
	HTTPDelete  HTTPMethod = "DELETE"  // RFC 7231, 4.3.5
	HTTPConnect HTTPMethod = "CONNECT" // RFC 7231, 4.3.6
	HTTPOptions HTTPMethod = "OPTIONS" // RFC 7231, 4.3.7
	HTTPTrace   HTTPMethod = "TRACE"   // RFC 7231, 4.3.8
	HTTPUse     HTTPMethod = "USE"
)

type EndpointResult struct {
	Code int
	Body reflect.Kind
}
type HandlerFunction func(c *fiber.Ctx, context *meta.ApplicationContext) error
type Endpoint struct {
	UrlExtension string
	AuthRequired AuthMethod
	Method       HTTPMethod
	Handler      HandlerFunction
	ResultSet    []EndpointResult
}
type EndpointOptions struct {
	AuthRequired AuthMethod
	Method       HTTPMethod
	ResultSet    []EndpointResult
}

type ApiDefinition struct {
	fiberApp      *fiber.App
	endpoints     []Endpoint
	Middleware    *MiddlewareConfiguration
	DDH           string
	AuthTokenName string
}

var DEFAULT_ENDPOINT_OPTIONS = EndpointOptions{
	AuthRequired: AuthDefault,
	Method:       HTTPGet,
	ResultSet:    []EndpointResult{},
}

func (apiDef *ApiDefinition) Add(url string, handler HandlerFunction, options *EndpointOptions) {
	if options == nil {
		options = &DEFAULT_ENDPOINT_OPTIONS
	}

	endpoint := Endpoint{
		UrlExtension: url,
		AuthRequired: meta.Ternary(options.AuthRequired == "", DEFAULT_ENDPOINT_OPTIONS.AuthRequired, options.AuthRequired),
		Method:       meta.Ternary(options.Method == "", DEFAULT_ENDPOINT_OPTIONS.Method, options.Method),
		Handler:      handler,
		ResultSet:    meta.Ternary(options.ResultSet == nil, DEFAULT_ENDPOINT_OPTIONS.ResultSet, options.ResultSet),
	}

	apiDef.endpoints = append(apiDef.endpoints, endpoint)
}

func New(fiberApp *fiber.App) (*ApiDefinition, error) {
	authTokenName, err := config.LoudGet("AUTH_TOKEN_NAME")
	if err != nil {
		return nil, err
	}
	return &ApiDefinition{
		fiberApp:      fiberApp,
		endpoints:     []Endpoint{},
		Middleware:    NewMiddlewareConfiguration(),
		DDH:           config.GetOr("DEFAULT_DEBUG_HEADER", "DEFAULT-DEBUG-HEADER"),
		AuthTokenName: authTokenName,
	}, nil
}

func (apiDef *ApiDefinition) BuildApi(context *meta.ApplicationContext) error {
	for _, handler := range apiDef.Middleware.GetAny() {
		apiDef.fiberApp.Use(handler)
	}
	log.Println("[api] Adding ", len(apiDef.endpoints), " endpoints")
	for _, endpoint := range apiDef.endpoints {
		var authMethodForEndpoint = getAuthForEndpoint(endpoint.AuthRequired)

		apiDef.fiberApp.Add(endpoint.Method, endpoint.UrlExtension, func(c *fiber.Ctx) error {
			authErr := authMethodForEndpoint(c, context)
			var risingEdgeMiddlewareErr error
			for _, handler := range apiDef.Middleware.GetRisingEdge() {
				if err := handler(c, context); err != nil {
					risingEdgeMiddlewareErr = err
				}
			}
			var endpointErr error
			if risingEdgeMiddlewareErr == nil && authErr == nil {
				endpointErr = endpoint.Handler(c, context)
			}
			for _, handler := range apiDef.Middleware.GetFallingEdge() {
				if err := handler(c, context); err != nil {
					return err
				}
			}
			return endpointErr
		})
		log.Println("[api]", endpoint.Method, endpoint.UrlExtension, " auth: ", endpoint.AuthRequired)
	}
	return nil
}

func getAuthForEndpoint(authType AuthMethod) func(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	switch authType {
	case AuthNone:
		return auth.NoAuth
	case AuthDefault:
		return func(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
			return auth.naiveCheckForHeaderAuth(c, "DEFAULT-DEBUG-HEADER")
		}
	}
	return nil
}

type MiddlewareConfiguration struct {
	risingEdge  []HandlerFunction
	fallingEdge []HandlerFunction
	any         []func(*fiber.Ctx) error
}

func NewMiddlewareConfiguration() *MiddlewareConfiguration {
	return &MiddlewareConfiguration{
		risingEdge:  make([]HandlerFunction, 10),
		fallingEdge: make([]HandlerFunction, 10),
		any:         []func(*fiber.Ctx) error{},
	}
}

func (midConf *MiddlewareConfiguration) AddFallingEdge(handler HandlerFunction) {
	midConf.fallingEdge = append(midConf.fallingEdge, handler)
}
func (midConf *MiddlewareConfiguration) GetRisingEdge() []HandlerFunction {
	return midConf.risingEdge
}
func (midConf *MiddlewareConfiguration) AddRisingEdge(handler HandlerFunction) {
	midConf.risingEdge = append(midConf.risingEdge, handler)
}
func (midConf *MiddlewareConfiguration) GetFallingEdge() []HandlerFunction {
	return midConf.fallingEdge
}
func (midConf *MiddlewareConfiguration) Add(fun func(*fiber.Ctx) error) {
	midConf.any = append(midConf.any, fun)
}
func (midConf *MiddlewareConfiguration) GetAny() []func(*fiber.Ctx) error {
	return midConf.any
}
