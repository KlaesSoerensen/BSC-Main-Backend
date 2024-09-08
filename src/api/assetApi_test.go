package api

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"otte_main_backend/src/meta"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Utility function to create and return a mocked GORM database
func createAssetGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Errorf("failed to open sqlmock database: %v", err)
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Errorf("failed to initialize gorm DB: %v", err)
		return nil, nil, err
	}

	return gormDB, mock, nil
}

// Utility function to create GORM mock and Fiber app
func setupAssetTest(t *testing.T) (gormDB *gorm.DB, mock sqlmock.Sqlmock, app *fiber.App, appContext *meta.ApplicationContext, err error) {
	gormDB, mock, err = createAssetGormMock(t)
	if err != nil {
		t.Errorf("Setup failed: %v", err)
		return nil, nil, nil, nil, err
	}

	app = fiber.New()
	appContext = meta.CreateApplicationContext(gormDB, gormDB, gormDB, "Test-DDH")

	err = applyAssetApi(app, appContext)
	if err != nil {
		t.Errorf("failed to apply asset API: %v", err)
		return nil, nil, nil, nil, err
	}

	return gormDB, mock, app, appContext, nil
}

// Utility function to test a request and validate the status code
func testAssetRequest(t *testing.T, app *fiber.App, method, path string, expectedStatusCode int) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("OTTE-Token", "OTTE-Token") // Set OTTE-Token for authentication

	resp, err := app.Test(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != expectedStatusCode {
		return resp, errors.New("unexpected status code")
	}

	return resp, nil
}

// Utility function to check mock expectations
func checkAssetMockExpectations(t *testing.T, mock sqlmock.Sqlmock) error {
	if err := mock.ExpectationsWereMet(); err != nil {
		return err
	}
	return nil
}

// Debugging method that logs useful information from gormDB and appContext, called only on test failure
func handleAppContextDebug(t *testing.T, gormDB *gorm.DB, appContext *meta.ApplicationContext) {
	// Log the DDH (Default Debug Header) for debugging purposes
	t.Logf("Debug Header (DDH): %s", appContext.DDH)

	// Log database connection statistics
	sqlDB, err := gormDB.DB()
	if err != nil {
		t.Errorf("Failed to get database connection from gormDB: %v", err)
	} else {
		stats := sqlDB.Stats()
		t.Logf("DB Stats: MaxOpenConnections=%d, OpenConnections=%d, InUseConnections=%d, IdleConnections=%d",
			stats.MaxOpenConnections, stats.OpenConnections, stats.InUse, stats.Idle)
	}

	// Log AuthTokenName for debugging
	t.Logf("Auth Token Name: %s", appContext.AuthTokenName)

	// Log ColonyAssetDB, LanguageDB, and PlayerDB details
	if appContext.ColonyAssetDB.NowFunc != nil {
		t.Logf("ColonyAssetDB.NowFunc: %v", appContext.ColonyAssetDB.NowFunc())
	}
	if appContext.LanguageDB.NowFunc != nil {
		t.Logf("LanguageDB.NowFunc: %v", appContext.LanguageDB.NowFunc())
	}
	t.Logf("PlayerDB: %+v", appContext.PlayerDB)
}

// Test function for getting an asset by ID
func TestGetAssetByID(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Mocking the expected row
	rows := sqlmock.NewRows([]string{"id", "alias", "type", "width", "height", "hasLODs", "id", "detailLevel", "blob"}).
		AddRow(1, "Test Asset", "image/png", 100, 100, false, 1, "high", nil)

	// Use case-insensitive matching for SQL query
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id = \$1 
		ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnRows(rows)

	// Make the request and check for errors
	resp, err := testAssetRequest(t, app, "GET", "/api/v1/asset/1", 200)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if there is an error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	// Check mock expectations
	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}

func TestNonexistentItem(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Use case-insensitive matching for SQL query
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id = \$1 
		ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(999, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/asset/999", 404)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}

func TestEmptyDatabase(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Use case-insensitive matching for SQL query
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id = \$1 
		ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(gorm.ErrRecordNotFound)

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/asset/1", 404)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}

func TestDatabaseConnectionError(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Use case-insensitive matching for SQL query
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id = \$1 
		ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("connection error"))

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/asset/1", 500)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}

func TestUnexpectedServerError(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Use case-insensitive matching for SQL query
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id = \$1 
		ORDER BY "GraphicalAsset"."id" LIMIT \$2`).
		WithArgs(1, 1).
		WillReturnError(errors.New("unexpected error"))

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/asset/1", 500)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}

func TestInvalidMultipleAssetIDs(t *testing.T) {
	// Assign all five returned values, but ignore 'mock' if it's not used
	gormDB, _, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Log the appContext.DDH for debugging
	handleAppContextDebug(t, gormDB, appContext)

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/assets?ids=abc,1", 400)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}
}

func TestNonexistentMultipleAssets(t *testing.T) {
	gormDB, mock, app, appContext, err := setupAssetTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Log the appContext.DDH for debugging
	handleAppContextDebug(t, gormDB, appContext)

	// Use case-insensitive matching for SQL query, with flexibility in case for `on`/`ON`
	mock.ExpectQuery(`(?i)SELECT "GraphicalAsset".*, "LOD".id, "LOD"."detailLevel", "LOD".blob 
		FROM "GraphicalAsset" 
		LEFT JOIN "LOD" ON "LOD"."graphicalAsset" = "GraphicalAsset".id 
		WHERE "GraphicalAsset".id IN \(\$1,\$2\)`).
		WithArgs(999, 1000).
		WillReturnError(gorm.ErrRecordNotFound)

	resp, err := testAssetRequest(t, app, "GET", "/api/v1/assets?ids=999,1000", 404)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only on error
		t.Errorf("Failed to process the request: %v", err)
	}

	if resp != nil {
		t.Logf("Response: %v", resp.Status)
	}

	err = checkAssetMockExpectations(t, mock)
	if err != nil {
		handleAppContextDebug(t, gormDB, appContext) // Log only if expectations were not met
		t.Errorf("Mock expectations not met: %v", err)
	}
}
