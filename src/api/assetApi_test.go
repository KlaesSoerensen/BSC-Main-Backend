package api

import (
	"errors"
	"net/http/httptest"
	"otte_main_backend/src/meta"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetAssetByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %v", err)
	}

	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "alias", "type", "width", "height", "hasLODs"}).
		AddRow(1, "Test Asset", "image/png", 100, 100, false)

	mock.ExpectQuery(`^SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetMultipleAssetsByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to initialize gorm DB: %v", err)
	}

	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	rows := sqlmock.NewRows([]string{"id", "alias", "type", "width", "height", "hasLODs"}).
		AddRow(1, "Test Asset 1", "image/png", 100, 100, false).
		AddRow(2, "Test Asset 2", "image/jpeg", 200, 200, true)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id IN \(\$1,\$2\)`).
		WithArgs(1, 2).
		WillReturnRows(rows)

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/assets?ids=1,2", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestIncorrectPath(t *testing.T) {
	app := fiber.New()
	applyAssetApi(app, meta.ApplicationContext{})

	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestNonexistentItem(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	mock.ExpectQuery(`^SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestInvalidItemID(t *testing.T) {
	app := fiber.New()
	applyAssetApi(app, meta.ApplicationContext{})

	req := httptest.NewRequest("GET", "/api/v1/asset/abc", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "text/plain")
}

func TestEmptyDatabase(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	mock.ExpectQuery(`^SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestDatabaseConnectionError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	mock.ExpectQuery(`^SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("connection error"))

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestUnexpectedServerError(t *testing.T) {
	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	appContext := meta.ApplicationContext{ColonyAssetDB: gormDB}

	mock.ExpectQuery(`^SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("unexpected error"))

	app := fiber.New()
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 500, resp.StatusCode)
}

func TestAssetIDZero(t *testing.T) {
	app := fiber.New()
	appContext := meta.ApplicationContext{ColonyAssetDB: &gorm.DB{}} // Ensure DB is initialized
	applyAssetApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/asset/0", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestRequestMultipleAssetsWithInvalidIDs(t *testing.T) {
	app := fiber.New()
	applyAssetApi(app, meta.ApplicationContext{})

	req := httptest.NewRequest("GET", "/api/v1/assets?ids=1,-2,abc", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 400, resp.StatusCode)
}
