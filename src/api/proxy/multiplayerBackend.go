package proxy

import (
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"
	"otte_main_backend/src/multiplayer"

	"github.com/gofiber/fiber/v2"
)

// Proxying some calls to get around browser pre-flight checks (CORS on non TLS or self-signed TLS)
func ApplyProxyAPI(app *fiber.App, context *meta.ApplicationContext) error {
	app.Get("/proxy/v1/multiplayer/lobby/:id", func(c *fiber.Ctx) error {
		return getLobbyStateProxyHandler(c, context)
	})

	return nil
}

func getLobbyStateProxyHandler(c *fiber.Ctx, context *meta.ApplicationContext) error {
	lobbyID, err := c.ParamsInt("id")
	if err != nil {
		c.Response().Header.Set(context.DDH, "Invalid lobby ID: "+err.Error())
		c.Status(fiber.StatusBadRequest)
		middleware.LogRequests(c)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	resp, err := multiplayer.GetLobbyState(uint32(lobbyID), context)
	if err != nil {
		c.Response().Header.Set(context.DDH, "Failed to get lobby state: "+err.Error())
		c.Status(fiber.StatusBadGateway)
		middleware.LogRequests(c)
		return c.SendStatus(fiber.StatusBadGateway)
	}

	c.Status(fiber.StatusOK)
	middleware.LogRequests(c)
	return c.JSON(resp)
}
