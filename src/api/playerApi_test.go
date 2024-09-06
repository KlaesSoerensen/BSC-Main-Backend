package api

import (
	"net/http/httptest"
	"testing"

	"net/http"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Utility function to test a request and check response
func testPlayerRequest(t *testing.T, app *fiber.App, method, path string, expectedStatusCode int) (*http.Response, error) {
	req := httptest.NewRequest(method, path, nil)
	resp, err := app.Test(req)

	if err != nil {
		return nil, err
	}

	assert.Equal(t, expectedStatusCode, resp.StatusCode)
	return resp, nil
}

// Utility function to check mock expectations
func checkPlayerMockExpectations(t *testing.T, mock sqlmock.Sqlmock) {
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlayerByID(t *testing.T) {
	mock, app, _ := setupTest(t)

	// Mock player data
	rows := sqlmock.NewRows([]string{"id", "IGN", "sprite"}).
		AddRow(1, "TestPlayer", 100)

	mock.ExpectQuery(`SELECT id, "IGN", sprite FROM "Player" WHERE id = $1`).
		WithArgs(1).
		WillReturnRows(rows)

	// Mock achievement count for tutorial completed
	mock.ExpectQuery(`SELECT count\(\*\) FROM "Achievement" WHERE player = $1 AND title = $2`).
		WithArgs(1, "Tutorial Completed").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock achievements data
	achievementRows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)

	mock.ExpectQuery(`SELECT id FROM "Achievement" WHERE player = \$1`).
		WithArgs(1).
		WillReturnRows(achievementRows)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/1", 200)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
	checkPlayerMockExpectations(t, mock)
}

func TestGetPlayerNotFound(t *testing.T) {
	mock, app, _ := setupTest(t)

	// Mock empty result set for non-existent player
	mock.ExpectQuery(`SELECT id, "IGN", sprite FROM "Player" WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/999", 404)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkPlayerMockExpectations(t, mock)
}

func TestInvalidPlayerID(t *testing.T) {
	_, app, _ := setupTest(t)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/abc", 400)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
}

func TestGetPlayerPreferences(t *testing.T) {
	mock, app, _ := setupTest(t)

	// Mock preferences data
	preferenceRows := sqlmock.NewRows([]string{"id", "preferenceKey", "chosenValue", "availableValues"}).
		AddRow(1, "Language", "EN", pq.Array([]string{"EN", "DK", "NO"}))

	mock.ExpectQuery(`SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(1).
		WillReturnRows(preferenceRows)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/1/preferences", 200)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
	checkPlayerMockExpectations(t, mock)
}

func TestGetPlayerPreferencesNotFound(t *testing.T) {
	mock, app, _ := setupTest(t)

	// Mock empty result set for non-existent preferences
	mock.ExpectQuery(`SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/999/preferences", 404)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	checkPlayerMockExpectations(t, mock)
}

func TestInvalidPlayerIDPreferences(t *testing.T) {
	_, app, _ := setupTest(t)

	_, err := testPlayerRequest(t, app, "GET", "/api/v1/player/abc/preferences", 400)
	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}
}
