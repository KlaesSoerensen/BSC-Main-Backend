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
	OverwritingSettings string `json:"overwritingSettings"`
}

type MinigameDifficultyDTO struct {
	ID                  uint32 `json:"id"`
	Name                string `json:"name"`
	Icon                uint32 `json:"icon"`
	Description         string `json:"description"`
	OverwritingSettings string `json:"overwritingSettings" gorm:"column:overwritingSettings"`
}

func (mdDTO *MinigameDifficultyDTO) TableName() string {
	return "MiniGameDifficulty"
}

type MinigameInfoDTO struct {
	ID           uint32                  `json:"id"`
	Name         string                  `json:"name"`
	Icon         uint32                  `json:"icon"`
	Description  string                  `json:"description"`
	Settings     string                  `json:"settings"`
	Difficulties []MinigameDifficultyDTO `json:"difficulties"`
}

func applyMinigameApi(app *fiber.App, appContext *meta.ApplicationContext) error {
	log.Println("[Minigame API] Applying Minigame API")

	app.Get("/api/v1/minigames/:id", auth.PrefixOn(appContext, getMinigameInfoHandler))

	return nil
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
		Table("MiniGame").
		Preload("MiniGameDifficulty").
		Where(`"MiniGame".id = ?`, id).
		Find(&minigame).Error; err != nil {

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
