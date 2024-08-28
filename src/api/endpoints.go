package api

import (
	"otte_main_backend/src/meta"

	"github.com/gofiber/fiber/v2"
)

func ApplyEndpoints(app *fiber.App, appContext meta.ApplicationContext) error {
	if err := applyCatalog(app, appContext); err != nil {
		return err
	}
	if err := applyAssetApi(app, appContext); err != nil {
		return err
	}
	if err := applyColonyApi(app, appContext); err != nil {
		return err
	}
	if err := applyCollectionApi(app, appContext); err != nil {
		return err
	}
	if err := applyHealthApi(app, appContext); err != nil {
		return err
	}
	if err := applyLocationApi(app, appContext); err != nil {
		return err
	}
	if err := applyLodApi(app, appContext); err != nil {
		return err
	}
	if err := applyMinigameApi(app, appContext); err != nil {
		return err
	}
	if err := applyMinigameApi(app, appContext); err != nil {
		return err
	}
	if err := applyPlayerApi(app, appContext); err != nil {
		return err
	}
	if err := applySessionApi(app, appContext); err != nil {
		return err
	}
	return nil
}
