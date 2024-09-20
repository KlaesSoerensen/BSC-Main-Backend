package api

import (
	"errors"
	"log"
	"otte_main_backend/src/meta"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// applyColonyApi sets up the colony API routes
func applyColonyApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Colony API] Applying colony API")

	// Route for opening a colony
	app.Post("/api/v1/colony/:colonyId/open", openColonyHandler(appContext))

	// Route for joining a colony
	app.Post("/api/v1/colony/join/:code", joinColonyHandler(appContext))

	return nil
}

// OpenColonyRequest represents the request body for opening a colony
type OpenColonyRequest struct {
	PlayerID uint32 `json:"playerId"`
}

// OpenColonyResponse represents the response for opening a colony
type OpenColonyResponse struct {
	Code                     string `json:"code"`
	LobbyID                  uint32 `json:"lobbyId"`
	MultiplayerServerAddress string `json:"multiplayerServerAddress"`
}

// JoinColonyResponse represents the response for joining a colony
type JoinColonyResponse struct {
	LobbyID                  uint32 `json:"lobbyId"`
	MultiplayerServerAddress string `json:"multiplayerServerAddress"`
}

// ColonyApiModel represents the Colony table for Colony API operations
type ColonyApiModel struct {
	ID         uint32             `gorm:"column:id;primaryKey"`
	Name       string             `gorm:"column:name"`
	AccLevel   int                `gorm:"column:accLevel"`
	Owner      uint32             `gorm:"column:owner"`
	ColonyCode ColonyCodeApiModel `gorm:"foreignKey:ColonyID"`
}

func (ColonyApiModel) TableName() string {
	return "Colony"
}

// ColonyCodeApiModel represents the ColonyCode table for Colony API operations
type ColonyCodeApiModel struct {
	ID            uint32 `gorm:"column:id;primaryKey"`
	LobbyID       uint32 `gorm:"column:lobbyId"`
	ServerAddress string `gorm:"column:serverAddress"`
	ColonyID      uint32 `gorm:"column:colony"`
	Value         string `gorm:"column:value"`
}

func (ColonyCodeApiModel) TableName() string {
	return "ColonyCode"
}

// openColonyHandler handles the request to open a colony
func openColonyHandler(appContext *meta.ApplicationContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		colonyID, err := c.ParamsInt("colonyId")
		if err != nil {
			c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
		}

		var req OpenColonyRequest
		if err := c.BodyParser(&req); err != nil {
			c.Response().Header.Set(appContext.DDH, "Invalid request body "+err.Error())
			return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
		}

		var colony ColonyApiModel
		if err := appContext.ColonyAssetDB.
			Preload("ColonyCode").
			Where("id = ? AND owner = ?", colonyID, req.PlayerID).
			First(&colony).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.Response().Header.Set(appContext.DDH, "Colony not found or not owned by player "+err.Error())
				return fiber.NewError(fiber.StatusNotFound, "Colony not found or not owned by player")
			}
			c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		response := OpenColonyResponse{
			Code:                     colony.ColonyCode.Value, // Use the Value field instead of ID
			LobbyID:                  colony.ColonyCode.LobbyID,
			MultiplayerServerAddress: colony.ColonyCode.ServerAddress,
		}

		return c.JSON(response)
	}
}

// joinColonyHandler handles the request to join a colony
func joinColonyHandler(appContext *meta.ApplicationContext) fiber.Handler {
	return func(c *fiber.Ctx) error {
		code := c.Params("code")

		// Check if the code is empty
		if code == "" {
			errorMsg := "Invalid colony code: Code is empty"
			log.Println(errorMsg)
			c.Response().Header.Set(appContext.DDH, errorMsg)
			return fiber.NewError(fiber.StatusBadRequest, errorMsg)
		}

		// Check if the code contains exactly six numeric characters
		matched, err := regexp.MatchString(`^\d{6}$`, code)
		if err != nil {
			errorMsg := "Error in code validation: " + err.Error()
			log.Println(errorMsg)
			c.Response().Header.Set(appContext.DDH, errorMsg)
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}
		if !matched {
			errorMsg := "Invalid colony code format: Code must be exactly 6 numeric characters"
			log.Println(errorMsg)
			c.Response().Header.Set(appContext.DDH, errorMsg)
			return fiber.NewError(fiber.StatusBadRequest, errorMsg)
		}

		var colonyCode ColonyCodeApiModel
		if err := appContext.ColonyAssetDB.
			Where("value = ?", code).
			First(&colonyCode).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				errorMsg := "Colony code not found: " + code
				log.Println(errorMsg)
				c.Response().Header.Set(appContext.DDH, errorMsg)
				return fiber.NewError(fiber.StatusNotFound, "Colony code not found")
			}
			errorMsg := "Database error: " + err.Error()
			log.Println(errorMsg)
			c.Response().Header.Set(appContext.DDH, errorMsg)
			return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
		}

		response := JoinColonyResponse{
			LobbyID:                  colonyCode.LobbyID,
			MultiplayerServerAddress: colonyCode.ServerAddress,
		}

		log.Printf("Successfully joined colony with code: %s", code)
		return c.JSON(response)
	}
}
