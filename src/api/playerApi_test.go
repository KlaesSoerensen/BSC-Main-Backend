package api

import (
	"net/http/httptest"
	"otte_main_backend/src/meta"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestGetPlayerByID(t *testing.T) {
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

	appContext := meta.ApplicationContext{PlayerDB: gormDB}

	// Mock player data
	rows := sqlmock.NewRows([]string{"id", "IGN", "sprite"}).
		AddRow(1, "TestPlayer", 100)

	mock.ExpectQuery(`^SELECT id, "IGN", sprite FROM "Player" WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	// Mock achievement count for tutorial completed
	mock.ExpectQuery(`^SELECT count\(\*\) FROM "Achievement" WHERE player = \$1 AND title = \$2`).
		WithArgs(1, "Tutorial Completed").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	// Mock achievements data
	achievementRows := sqlmock.NewRows([]string{"id"}).
		AddRow(1).
		AddRow(2)

	mock.ExpectQuery(`^SELECT id FROM "Achievement" WHERE player = \$1`).
		WithArgs(1).
		WillReturnRows(achievementRows)

	app := fiber.New()
	applyPlayerApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/player/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlayerNotFound(t *testing.T) {
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

	appContext := meta.ApplicationContext{PlayerDB: gormDB}

	// Mock empty result set for non-existent player
	mock.ExpectQuery(`^SELECT id, "IGN", sprite FROM "Player" WHERE id = \$1`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	app := fiber.New()
	applyPlayerApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/player/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestInvalidPlayerID(t *testing.T) {
	app := fiber.New()
	applyPlayerApi(app, meta.ApplicationContext{})

	req := httptest.NewRequest("GET", "/api/v1/player/abc", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
}

func TestGetPlayerPreferences(t *testing.T) {
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

	appContext := meta.ApplicationContext{PlayerDB: gormDB}

	// Mock preferences data
	preferenceRows := sqlmock.NewRows([]string{"id", "preferenceKey", "chosenValue", "availableValues"}).
		AddRow(1, "Language", "EN", pq.Array([]string{"EN", "DK", "NO"}))

	mock.ExpectQuery(`^SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(1).
		WillReturnRows(preferenceRows)

	app := fiber.New()
	applyPlayerApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/player/1/preferences", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 200, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetPlayerPreferencesNotFound(t *testing.T) {
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

	appContext := meta.ApplicationContext{PlayerDB: gormDB}

	// Mock empty result set for non-existent preferences
	mock.ExpectQuery(`^SELECT "PlayerPreference".id, "PlayerPreference"."preferenceKey", "PlayerPreference"."chosenValue", "AvailablePreference"."availableValues" FROM "PlayerPreference" JOIN "AvailablePreference"`).
		WithArgs(999).
		WillReturnError(gorm.ErrRecordNotFound)

	app := fiber.New()
	applyPlayerApi(app, appContext)

	req := httptest.NewRequest("GET", "/api/v1/player/999/preferences", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, 404, resp.StatusCode)
}

func TestInvalidPlayerIDPreferences(t *testing.T) {
	app := fiber.New()
	applyPlayerApi(app, meta.ApplicationContext{})

	req := httptest.NewRequest("GET", "/api/v1/player/abc/preferences", nil)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("Failed to process the request: %v", err)
	}

	assert.Equal(t, 400, resp.StatusCode)
	assert.Contains(t, resp.Header.Get("Content-Type"), "application/json")
}
