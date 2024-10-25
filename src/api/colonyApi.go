package api

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"otte_main_backend/src/middleware"
	"otte_main_backend/src/multiplayer"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func applyColonyApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	app.Get("/api/v1/colony/:colonyId/pathgraph", auth.PrefixOn(appContext, getPathGraphHandler))
	app.Get("/api/v1/colony/:colonyId/code", auth.PrefixOn(appContext, getColonyCodeHandler))
	app.Post("/api/v1/colony/:colonyId/open", auth.PrefixOn(appContext, openColonyHandler))
	app.Post("/api/v1/colony/:colonyId/close", func(c *fiber.Ctx) error {
		return closeColonyHandler(c, appContext)
	})
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

type OpenColonyRequest struct {
	DurationMS  uint32 `json:"validDurationMS"` // Duration in milliseconds
	PlayerID    uint32 `json:"playerId"`
	LatestVisit string `json:"latestVisit"`
}

type OpenColonyResponse struct {
	Code                     string `json:"code"`
	LobbyID                  uint32 `json:"lobbyId"`
	MultiplayerServerAddress string `json:"multiplayerServerAddress"`
}

type JoinColonyResponse struct {
	LobbyID                  uint32 `json:"lobbyId"`
	MultiplayerServerAddress string `json:"multiplayerServerAddress"`
	Owner                    uint32 `json:"owner"`
	ColonyID                 uint32 `json:"colonyId"`
}

type ColonyDTO struct {
	ID          uint32           `gorm:"column:id;primaryKey"`
	Name        string           `gorm:"column:name"`
	AccLevel    int              `gorm:"column:accLevel"`
	Owner       uint32           `gorm:"column:owner"`
	LatestVisit string           `gorm:"column:latestVisit"`
	ColonyCode  *ColonyCodeModel `gorm:"foreignKey:ColonyID"`
}

func (ColonyDTO) TableName() string {
	return "Colony"
}

type ColonyCodeModel struct {
	ID              uint32    `gorm:"column:id;primaryKey"`
	LobbyID         uint32    `gorm:"column:lobbyId"`
	ServerAddress   string    `gorm:"column:serverAddress"`
	ColonyID        uint32    `gorm:"column:colony"`
	Value           string    `gorm:"column:value"`
	OwnerID         uint32    `gorm:"column:owner"`
	CreatedAt       time.Time `gorm:"column:createdAt"`
	ValidDurationMS uint32    `gorm:"column:validDurationMS"`
}

func (ColonyCodeModel) TableName() string {
	return "ColonyCode"
}

func openColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	colonyID, err := c.ParamsInt("colonyId")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	var req OpenColonyRequest
	if err := c.BodyParser(&req); err != nil || req.PlayerID == 0 {
		c.Response().Header.Set(appContext.DDH, "Invalid request body "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.DurationMS == 0 {
		req.DurationMS = 600000
	}

	var colony ColonyDTO
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

	if colony.ColonyCode != nil {
		if colony.ColonyCode.CreatedAt.Add(time.Duration(colony.ColonyCode.ValidDurationMS) * time.Millisecond).After(time.Now()) {
			response := OpenColonyResponse{
				Code:                     colony.ColonyCode.Value,
				LobbyID:                  colony.ColonyCode.LobbyID,
				MultiplayerServerAddress: colony.ColonyCode.ServerAddress,
			}
			c.Status(fiber.StatusOK)
			return c.JSON(response)
		}
	}

	lobbyID, err := multiplayer.CreateLobby(req.PlayerID, colony.ID, appContext)
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Failed to create lobby "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to create lobby")
	}

	colony.LatestVisit = req.LatestVisit
	colony.ColonyCode = &ColonyCodeModel{
		LobbyID:         lobbyID,
		ServerAddress:   appContext.ExternalMultiplayerServerAddress,
		ColonyID:        colony.ID,
		OwnerID:         req.PlayerID,
		ValidDurationMS: req.DurationMS,
		CreatedAt:       time.Now(),
	}

	var isGood = false
	var retryCount = 0
	const maxRetries = 10

	for !isGood && retryCount < maxRetries {
		maxNum := big.NewInt(1000000)
		n, err := rand.Int(rand.Reader, maxNum)
		if err != nil {
			c.Response().Header.Set(appContext.DDH, "Failed to generate secure random number "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate colony code")
		}
		asString := fmt.Sprintf("%06d", n.Int64())
		backToInt, _ := strconv.Atoi(asString)
		if backToInt < 100000 { // limiting range to 100000-999999
			backToInt += 100000
		}

		colony.ColonyCode.Value = fmt.Sprintf("%d", backToInt)

		tx := appContext.ColonyAssetDB.Begin()
		if err := tx.Error; err != nil {
			c.Response().Header.Set(appContext.DDH, "Failed to begin transaction "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to begin transaction")
		}

		if err := tx.Create(colony.ColonyCode).Error; err != nil {
			tx.Rollback()
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				retryCount++
				continue
			}
			c.Response().Header.Set(appContext.DDH, "Failed to create ColonyCode "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to create ColonyCode")
		}

		if err := tx.Model(&colony).Update("colonyCode", colony.ColonyCode.ID).Error; err != nil {
			tx.Rollback()
			c.Response().Header.Set(appContext.DDH, "Failed to update Colony colonyCode "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to update Colony")
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			c.Response().Header.Set(appContext.DDH, "Failed to commit transaction "+err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
		}
		isGood = true
	}

	if retryCount >= maxRetries {
		c.Response().Header.Set(appContext.DDH, "Failed to generate unique colony code after maximum retries")
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to generate unique colony code")
	}

	response := OpenColonyResponse{
		Code:                     colony.ColonyCode.Value,
		LobbyID:                  colony.ColonyCode.LobbyID,
		MultiplayerServerAddress: colony.ColonyCode.ServerAddress,
	}
	c.Status(fiber.StatusOK)
	return c.JSON(response)
}

func joinColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	code := c.Params("code")

	if code == "" {
		c.Response().Header.Set(appContext.DDH, "Invalid colony code: Code is empty")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony code: Code is empty")
	}

	matched, err := regexp.MatchString(`^\d{6}$`, code)
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Error in code validation: "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}
	if !matched {
		c.Response().Header.Set(appContext.DDH, "Invalid colony code format: Code must be exactly 6 numeric characters")
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony code format")
	}

	var colonyCodes []ColonyCodeModel
	if err := appContext.ColonyAssetDB.
		Where("value = ?", code).
		Find(&colonyCodes).Error; err != nil {
		c.Response().Header.Set(appContext.DDH, "Database error: "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if len(colonyCodes) == 0 {
		c.Response().Header.Set(appContext.DDH, "Colony code not found: "+code)
		return fiber.NewError(fiber.StatusNotFound, "Colony code not found")
	}

	colonyCode := colonyCodes[0]
	if colonyCode.CreatedAt.Add(time.Duration(colonyCode.ValidDurationMS) * time.Millisecond).Before(time.Now()) {
		appContext.ColonyAssetDB.Delete(&colonyCode)
		c.Response().Header.Set(appContext.DDH, "Code expired")
		return fiber.NewError(fiber.StatusNotFound, "Colony code not found: "+code)
	}

	response := JoinColonyResponse{
		Owner:                    colonyCode.OwnerID,
		LobbyID:                  colonyCode.LobbyID,
		MultiplayerServerAddress: colonyCode.ServerAddress,
		ColonyID:                 colonyCode.ColonyID,
	}
	c.Status(fiber.StatusOK)
	return c.JSON(response)
}

type UpdateLatestVisitRequest struct {
	LatestVisit string `json:"latestVisit"`
}

type UpdateLatestVisitResponse struct {
	LatestVisit string `json:"latestVisit"`
}

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

	var colony ColonyDTO
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

type CloseColonyRequest struct {
	PlayerID uint32 `json:"playerId"`
}

func closeColonyHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	colonyID, err := c.ParamsInt("colonyId")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
		c.Status(fiber.StatusBadRequest)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	var req CloseColonyRequest
	if err := c.BodyParser(&req); err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid request body "+err.Error())
		c.Status(fiber.StatusBadRequest)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}
	if req.PlayerID == 0 {
		c.Response().Header.Set(appContext.DDH, "Invalid request body: Player ID is 0")
		c.Status(fiber.StatusBadRequest)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tx := appContext.ColonyAssetDB.Begin()
	if err := tx.Error; err != nil {
		c.Response().Header.Set(appContext.DDH, "Failed to begin transaction "+err.Error())
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to begin transaction")
	}

	var colony ColonyDTO
	if err := tx.Where("id = ? AND owner = ?", colonyID, req.PlayerID).First(&colony).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony not found or not owned by player "+err.Error())
			c.Status(fiber.StatusNotFound)
			middleware.LogRequests(c)
			return fiber.NewError(fiber.StatusNotFound, "Colony not found or not owned by player")
		}
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if err := tx.Model(&colony).Update("colonyCode", nil).Error; err != nil {
		tx.Rollback()
		c.Response().Header.Set(appContext.DDH, "Failed to update Colony colonyCode "+err.Error())
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to update Colony")
	}

	if err := tx.Where("colony = ?", colonyID).Delete(&ColonyCodeModel{}).Error; err != nil {
		tx.Rollback()
		c.Response().Header.Set(appContext.DDH, "Failed to delete ColonyCode "+err.Error())
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to delete ColonyCode")
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.Response().Header.Set(appContext.DDH, "Failed to commit transaction "+err.Error())
		c.Status(fiber.StatusInternalServerError)
		middleware.LogRequests(c)
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to commit transaction")
	}

	c.Status(fiber.StatusOK)
	middleware.LogRequests(c)
	return c.SendStatus(fiber.StatusOK)
}

func getColonyCodeHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	colonyID, err := c.ParamsInt("colonyId")
	if err != nil {
		c.Response().Header.Set(appContext.DDH, "Invalid colony ID "+err.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Invalid colony ID")
	}

	// First get the colony to find its colonyCode ID
	var colony struct {
		ID         uint32 `gorm:"column:id"`
		ColonyCode uint32 `gorm:"column:colonyCode"`
	}
	if err := appContext.ColonyAssetDB.
		Table("Colony").
		Select("id, \"colonyCode\"").
		Where("id = ?", colonyID).
		First(&colony).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Colony not found")
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	if colony.ColonyCode == 0 {
		c.Response().Header.Set(appContext.DDH, "No active colony code found")
		return fiber.NewError(fiber.StatusNotFound, "No active colony code found")
	}

	// Now get the actual colony code using the ID from Colony table
	var colonyCode ColonyCodeModel
	if err := appContext.ColonyAssetDB.
		Table("ColonyCode").
		Where("id = ?", colony.ColonyCode).
		First(&colonyCode).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "Colony code not found "+err.Error())
			return fiber.NewError(fiber.StatusNotFound, "Colony code not found")
		}
		c.Response().Header.Set(appContext.DDH, "Internal server error "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Internal server error")
	}

	// Check if the code has expired
	if colonyCode.CreatedAt.Add(time.Duration(colonyCode.ValidDurationMS) * time.Millisecond).Before(time.Now()) {
		c.Response().Header.Set(appContext.DDH, "Colony code has expired")
		return fiber.NewError(fiber.StatusNotFound, "Colony code has expired")
	}

	response := OpenColonyResponse{
		Code: colonyCode.Value,
	}

	c.Status(fiber.StatusOK)
	return c.JSON(response)
}
