package api

import (
	"errors"
	"net/http/httptest"
	"otte_main_backend/src/meta"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Utility function to create and return a mocked GORM database
func createPlayerGormMock(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Error("failed to open sqlmock database:", err)
		return nil, nil, err
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		t.Error("failed to initialize gorm DB:", err)
		return nil, nil, err
	}

	return gormDB, mock, nil
}

// Utility function to create GORM mock and Fiber app
func setupPlayerTest(t *testing.T) (sqlmock.Sqlmock, *fiber.App, *meta.ApplicationContext, error) {
	gormDB, mock, err := createPlayerGormMock(t)
	if err != nil {
		t.Error("Setup failed:", err)
		return nil, nil, nil, err
	}

	app := fiber.New()
	appContext := &meta.ApplicationContext{
		ColonyAssetDB:            gormDB,
		LanguageDB:               gormDB,
		PlayerDB:                 gormDB,
		VitecIntegration:         nil,
		DDH:                      "Test-DDH",
		AuthTokenName:            "Test-Auth-Token-Name",
		MultiplayerServerAddress: "notset",
	}

	err = applyPlayerApi(app, appContext)
	if err != nil {
		t.Error("failed to apply player API:", err)
		return nil, nil, nil, err
	}

	return mock, app, appContext, nil
}

// Utility function to test a request and validate the status code
func testPlayerRequest(t *testing.T, app *fiber.App, method, path string, expectedStatusCode int) error {
	req := httptest.NewRequest(method, path, nil)
	req.Header.Set("OTTE-Token", "OTTE-Token") // Set OTTE-Token for authentication

	resp, err := app.Test(req)
	if err != nil {
		t.Error("failed to process the request:", err)
		return err
	}

	if resp.StatusCode != expectedStatusCode {
		t.Errorf("unexpected status code: got %d, expected %d", resp.StatusCode, expectedStatusCode)
		return errors.New("unexpected status code")
	}

	return nil
}

// Utility function to check mock expectations
func checkPlayerMockExpectations(t *testing.T, mock sqlmock.Sqlmock) error {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Error("there were unfulfilled expectations:", err)
		return err
	}
	return nil
}

func TestGetPlayerByID(t *testing.T) {
	mock, app, _, err := setupPlayerTest(t) // Capture all 4 return values
	if err != nil {
		t.Fatal("Setup failed:", err) // Fatal here as setup failure should stop the test
	}

	// Mocking the expected row
	rows := sqlmock.NewRows([]string{"id", "IGN", "sprite"}).
		AddRow(1, "TestPlayer", 100)

	mock.ExpectQuery(`SELECT id, "IGN", sprite FROM "Player" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/1", 200)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}

	err = checkPlayerMockExpectations(t, mock)
	if err != nil {
		t.Error("Mock expectations not met:", err)
	}
}

func TestGetPlayerNotFound(t *testing.T) {
	mock, app, _, err := setupPlayerTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	mock.ExpectQuery(`SELECT id, "IGN", sprite FROM "Player" WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/999", 404)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}

	err = checkPlayerMockExpectations(t, mock)
	if err != nil {
		t.Error("Mock expectations not met:", err)
	}
}

func TestInvalidPlayerID(t *testing.T) {
	_, app, _, err := setupPlayerTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/abc", 400)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}
}

func TestGetPlayerPreferences(t *testing.T) {
	mock, app, _, err := setupPlayerTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	// Mocking the preferences row
	preferenceRows := sqlmock.NewRows([]string{"id", "preferenceKey", "chosenValue", "availableValues"}).
		AddRow(1, "Language", "EN", pq.Array([]string{"EN", "DK", "NO"}))

	mock.ExpectQuery(`SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(1).
		WillReturnRows(preferenceRows)

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/1/preferences", 200)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}

	err = checkPlayerMockExpectations(t, mock)
	if err != nil {
		t.Error("Mock expectations not met:", err)
	}
}

func TestGetPlayerPreferencesNotFound(t *testing.T) {
	mock, app, _, err := setupPlayerTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	mock.ExpectQuery(`SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/999/preferences", 404)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}

	err = checkPlayerMockExpectations(t, mock)
	if err != nil {
		t.Error("Mock expectations not met:", err)
	}
}

func TestInvalidPlayerIDPreferences(t *testing.T) {
	_, app, _, err := setupPlayerTest(t)
	if err != nil {
		t.Fatal("Setup failed:", err)
	}

	err = testPlayerRequest(t, app, "GET", "/api/v1/player/abc/preferences", 400)
	if err != nil {
		t.Error("Failed to process the request:", err)
	}
}
