package api

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/multiplayer"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// applyColonyApi sets up the colony API routes
func applyColonyApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Colony API] Applying colony API")

	app.Get("/api/v1/colony/:colonyId/pathgraph", auth.PrefixOn(appContext, getPathGraphHandler))
	app.Post("/api/v1/colony/:colonyId/open", auth.PrefixOn(appContext, openColonyHandler))
	app.Post("/api/v1/colony/join/:code", auth.PrefixOn(appContext, joinColonyHandler))
	app.Post("/api/v1/colony/:colonyId/update-last-visit", auth.PrefixOn(appContext, updateLatestVisitHandler))

	return nil
}

type PathDTO struct {
	From uint32 `json:"from" gorm:"column:locationA"` //Id of ColonyLocation
	To   uint32 `json:"to" gorm:"column:locationB"`   //Id of ColonyLocation
}

func (p *PathDTO) TableName() string {
	return "ColonyLocationPath"
}

type PathGraphDTO struct {
	Paths []PathDTO `json:"paths"`
}

func getPathGraphHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	colonyID, err := c.ParamsInt("colonyId")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}
	var paths []PathDTO
	if dbErr := appContext.ColonyAssetDB.Where("colony = ?", colonyID).Find(&paths).Error; dbErr != nil || len(paths) == 0 {
		if !errors.Is(dbErr, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Internal error")
			return fiber.NewError(fiber.StatusInternalServerError, "Internal error")
		}

		c.Response().Header.Set(appContext.DDH, "Colony not found or paths not found")
		return fiber.NewError(fiber.StatusNotFound, "Colony not found or paths not found")
	}
	c.Status(fiber.StatusOK)
	return c.JSON(PathGraphDTO{Paths: paths})
}

// OpenColonyRequest represents the request body for opening a colony
type OpenColonyRequest struct {
	PlayerID    uint32 `json:"playerId"`
	LatestVisit string `json:"latestVisit"`
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
	Owner                    uint32 `json:"owner"`
}

// ColonyApiModel represents the Colony table for Colony API operations
type ColonyApiModel struct {
	ID          uint32             `gorm:"column:id;primaryKey"`
	Name        string             `gorm:"column:name"`
	AccLevel    int                `gorm:"column:accLevel"`
	Owner       uint32             `gorm:"column:owner"`
	LatestVisit string             `gorm:"column:latestVisit"`
	ColonyCode  ColonyCodeApiModel `gorm:"foreignKey:ColonyID"`
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
	OwnerID       uint32 `gorm:"column:owner"`
}

func (ColonyCodeApiModel) TableName() string {
	return "ColonyCode"
}

// openColonyHandler handles the request to open a colony
func openColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
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

	lobbyID, err := multiplayer.CreateLobby(req.PlayerID, appContext)
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Failed to create lobby "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create lobby")
	}

	// Update the LatestVisit field with the provided value
	colony.LatestVisit = req.LatestVisit
	colony.ColonyCode.OwnerID = colony.Owner
	colony.ColonyCode.LobbyID = lobbyID
	colony.ColonyCode.ServerAddress = appContext.MultiplayerServerAddress

	var isGood = false
	for !isGood {
		code := generateColonyCode()
		colony.ColonyCode.Value = code

		if err := appContext.ColonyAssetDB.Save(&colony).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				isGood = false
				continue
			}

			c.Response().Header.Set(appContext.DDH, "Failed to update ColonyCode "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update ColonyCode")
		}
		isGood = true
	}

	response := OpenColonyResponse{
		Code:                     colony.ColonyCode.Value,
		LobbyID:                  colony.ColonyCode.LobbyID,
		MultiplayerServerAddress: colony.ColonyCode.ServerAddress,
	}

	return c.JSON(response)
}

var localRand = rand.New(rand.NewSource(0))

// Generates a code with 6 digits
func generateColonyCode() string {
	number := localRand.Intn(1000000)
	return fmt.Sprintf("%06d", number)
}

// joinColonyHandler handles the request to join a colony
func joinColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
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
		c.Response().Header.Set(appContext.DDH, errorMsg)
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	if !matched {
		errorMsg := "Invalid colony code format: Code must be exactly 6 numeric characters"
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
		Owner:                    colonyCode.OwnerID,
		LobbyID:                  colonyCode.LobbyID,
		MultiplayerServerAddress: colonyCode.ServerAddress,
	}

	log.Printf("Successfully joined colony with code: %s", code)
	return c.JSON(response)
}

// UpdateLatestVisitRequest represents the request body for updating the latest visit time
type UpdateLatestVisitRequest struct {
	LatestVisit string `json:"latestVisit"`
}

// UpdateLatestVisitResponse represents the response body after updating the latest visit time
type UpdateLatestVisitResponse struct {
	LatestVisit string `json:"latestVisit"`
}

// Add this new handler function
func updateLatestVisitHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	colonyID, err := c.ParamsInt("colonyId")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	var req UpdateLatestVisitRequest
	if err := c.BodyParser(&req); err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid request body "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	var colony ColonyApiModel
	if err := appContext.ColonyAssetDB.
		Where("id = ?", colonyID).
		First(&colony).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony not found or not owned by player "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Colony not found or not owned by player")
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Update the LatestVisit field with the provided value
	colony.LatestVisit = req.LatestVisit
	if err := appContext.ColonyAssetDB.Save(&colony).Error; err != nil {
		c.Response().Header.Set(appContext.DDH, "Failed to update LatestVisit "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update LatestVisit")
	}

	response := UpdateLatestVisitResponse{
		LatestVisit: colony.LatestVisit,
	}

	return c.JSON(response)
}
