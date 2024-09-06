package api

import (
	"errors"
	"net/http/httptest"
	"otte_main_backend/src/meta"
	"testing"

	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Utility function to create and return a mocked GORM database
func createGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
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

	return gormDB, mock
}

// Utility function to create GORM mock and Fiber app
func setupTest(t *testing.T) (sqlmock.Sqlmock, *fiber.App, *meta.ApplicationContext) {
	gormDB, mock := createGormMock(t)
	app := fiber.New()
	appContext := meta.CreateApplicationContext(gormDB, gormDB, gormDB, "Test-DDH")

	err := applyAssetApi(app, appContext)
	if err != nil {
		t.Fatalf("failed to apply asset API: %v", err)
	}

	return mock, app, appContext
}

// Utility function to test a request and check response
func testRequest(t *testing.T, app *fiber.App, method, path string, expectedStatusCode int) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req)

	if err != nil {
		return nil, err
	}

	assert.Equal(t, expectedStatusCode, resp.StatusCode)
	return resp, nil
}

// Utility function to check mock expectations
func checkMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAssetByID(t *testing.T) {
	mock, app, _ := setupTest(t)

	rows := sqlmock.NewRows([]string{"id", "alias", "type", "width", "height", "hasLODs"}).
		AddRow(1, "Test Asset", "image/png", 100, 100, false)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	_, err := testRequest(t, app, "GET", "/api/v1/asset/1", 200)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
	checkMockExpectations(t, mock)
}

func TestNonexistentItem(t *testing.T) {
	mock, app, _ := setupTest(t)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := testRequest(t, app, "GET", "/api/v1/asset/999", 404)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkMockExpectations(t, mock)
}

func TestEmptyDatabase(t *testing.T) {
	mock, app, _ := setupTest(t)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := testRequest(t, app, "GET", "/api/v1/asset/1", 404)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkMockExpectations(t, mock)
}

func TestDatabaseConnectionError(t *testing.T) {
	mock, app, _ := setupTest(t)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("connection error"))

	_, err := testRequest(t, app, "GET", "/api/v1/asset/1", 500)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkMockExpectations(t, mock)
}

func TestUnexpectedServerError(t *testing.T) {
	mock, app, _ := setupTest(t)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id = \$1 ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("unexpected error"))

	_, err := testRequest(t, app, "GET", "/api/v1/asset/1", 500)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkMockExpectations(t, mock)
}

func TestInvalidMultipleAssetIDs(t *testing.T) {
	_, app, _ := setupTest(t)

	_, err := testRequest(t, app, "GET", "/api/v1/assets?ids=abc,1", 400)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
}

func TestNonexistentMultipleAssets(t *testing.T) {
	mock, app, _ := setupTest(t)

	mock.ExpectQuery(`SELECT "GraphicalAsset".*, "LOD".id as lod_id, "LOD"."detailLevel" as detail_level, "LOD".blob as lod_blob FROM "GraphicalAsset" LEFT JOIN "LOD" on "LOD"."graphicalAsset" = "GraphicalAsset".id WHERE "GraphicalAsset".id IN \(\$1,\$2\)`).
		WithArgs(999, 1000).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := testRequest(t, app, "GET", "/api/v1/assets?ids=999,1000", 404)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkMockExpectations(t, mock)
}
