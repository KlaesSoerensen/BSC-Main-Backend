package api

import (
	"errors"
	"log"
	"otte_main_backend/src/auth"
	"otte_main_backend/src/meta"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type MinimizedMinigameDTO struct {
	Settings            string `json:"settings"`
	OverwritingSettings string `json:"overwritingSettings" gorm:"column:overwritingSettings"`
}

type MinigameDifficultyDTO struct {
	ID                  uint32 `json:"id"`
	Name                string `json:"name"`
	MinigameID          uint32 `json:"-" gorm:"column:minigame"`
	Icon                uint32 `json:"icon"`
	Description         string `json:"description"`
	OverwritingSettings string `json:"overwritingSettings" gorm:"column:overwritingSettings"`
}

func (mdDTO *MinigameDifficultyLocationDTO) TableName() string {
	return "MiniGameDifficulty"
}

type MinigameInfoDTO struct {
	ID           uint32                          `json:"id"`
	Name         string                          `json:"name"`
	Icon         uint32                          `json:"icon"`
	Description  string                          `json:"description"`
	Settings     string                          `json:"settings"`
	Difficulties []MinigameDifficultyLocationDTO `json:"difficulties" gorm:"foreignKey:MinigameID;references:ID"`
}

func (mInfoDTO *MinigameInfoDTO) TableName() string {
	return "MiniGame"
}

func applyMinigameApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Minigame API] Applying Minigame API")

	//Fiber is so stateful that it infers "minimized" as a query param if registered after minigames/:id
	//y
	app.Get("/api/v1/minigame/minimized", auth.PrefixOn(appContext, getMinimizedMinigameHandler))

	app.Get("/api/v1/minigame/:id", auth.PrefixOn(appContext, getMinigameInfoHandler))

	return nil
}

func getMinimizedMinigameHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	idStr := c.Query("minigame")
	diffStr := c.Query("difficulty")
	minigameID, idParseErr := strconv.Atoi(idStr)
	if idParseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Error in parsing minigame id "+idParseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Error in parsing minigame id "+idParseErr.Error())
	}
	diffID, diffParseErr := strconv.Atoi(diffStr)
	if diffParseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Error in parsing minigame difficulty "+diffParseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Error in parsing minigame difficulty "+diffParseErr.Error())
	}

	var minigame MinimizedMinigameDTO
	if err := appContext.ColonyAssetDB.
		Table("MiniGame").
		Table("MiniGame").
		Select(`"MiniGame".settings, "MiniGameDifficulty"."overwritingSettings"`).
		Joins(`JOIN "MiniGameDifficulty" ON "MiniGame".id = ?`, minigameID).
		Where(`"MiniGameDifficulty".id = ?`, diffID).
		Scan(&minigame).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No such minigame or difficulty")
			return fiber.NewError(fiber.StatusNotFound, "No such minigame or difficulty")
		}

		c.Response().Header.Set(appContext.DDH, "Error in fetching minigame "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error in fetching minigame "+err.Error())
	}

	if minigame.Settings == "" || minigame.OverwritingSettings == "" {
		c.Response().Header.Set(appContext.DDH, "No such minigame or minigame difficulty")
		return fiber.NewError(fiber.StatusNotFound, "No such minigame or minigame difficulty")
	}

	c.Status(fiber.StatusOK)
	return c.JSON(minigame)
}

func getMinigameInfoHandler(c *fiber.Ctx, appContext *meta.ApplicationContext) error {
	idStr := c.Params("id")
	id, parseErr := strconv.Atoi(idStr)

	if parseErr != nil {
		c.Response().Header.Set(appContext.DDH, "Error in parsing minigame id "+parseErr.Error())
		return fiber.NewError(fiber.StatusBadRequest, "Error in parsing minigame id "+parseErr.Error())
	}

	var minigame MinigameInfoDTO
	if err := appContext.ColonyAssetDB.
		Preload("Difficulties").
		Where(`"MiniGame".id = ?`, id).
		First(&minigame).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Response().Header.Set(appContext.DDH, "No such minigame")
			return fiber.NewError(fiber.StatusNotFound, "No such minigame")
		}

		c.Response().Header.Set(appContext.DDH, "Error in fetching minigame "+err.Error())
		return fiber.NewError(fiber.StatusInternalServerError, "Error in fetching minigame "+err.Error())
	}

	c.Status(fiber.StatusOK)
	return c.JSON(minigame)
}
