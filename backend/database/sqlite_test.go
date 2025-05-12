package database

import (
	"homelab-monitor/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to test database")
	}
	_ = DB.AutoMigrate(&models.User{}, &models.UserAuth{}, &models.SpeedTest{}, &models.DailyApiCall{})
}

func TestInitDB(t *testing.T) {
	setupTestDB()
	if DB == nil {
		t.Fatal("DB should not be nil after initialization")
	}
}

func TestGetAllUsers(t *testing.T) {
	setupTestDB()
	DB.Create(&models.User{Email: "test@example.com"})

	users, err := GetAllUsers()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestGetUserByEmail(t *testing.T) {
	setupTestDB()
	DB.Create(&models.User{Email: "test@example.com"})

	user, err := GetUserByEmail("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("expected email 'test@example.com', got '%s'", user.Email)
	}
}

func TestUpdateUserPhoto(t *testing.T) {
	setupTestDB()
	DB.Create(&models.User{Email: "test@example.com"})

	err := UpdateUserPhoto("test@example.com", "/path/to/photo.jpg")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	user, _ := GetUserByEmail("test@example.com")
	if user.PhotoURL != "/path/to/photo.jpg" {
		t.Fatalf("expected photo URL '/path/to/photo.jpg', got '%s'", user.PhotoURL)
	}
}

func TestGetUserPhotoPath(t *testing.T) {
	setupTestDB()
	DB.Create(&models.User{Email: "test@example.com", PhotoURL: "/path/to/photo.jpg"})

	photoPath, err := GetUserPhotoPath("test@example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if photoPath != "/path/to/photo.jpg" {
		t.Fatalf("expected photo path '/path/to/photo.jpg', got '%s'", photoPath)
	}
}

func TestGetDailyApiCalls(t *testing.T) {
	setupTestDB()
	now := time.Now().Truncate(24 * time.Hour)
	DB.Create(&models.DailyApiCall{Date: now.AddDate(0, 0, -1)})
	DB.Create(&models.DailyApiCall{Date: now.AddDate(0, 0, -15)})

	results, err := GetDailyApiCalls()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !results[0].Date.Equal(now.AddDate(0, 0, -1)) {
		t.Fatalf("expected date '%v', got '%v'", now.AddDate(0, 0, -1), results[0].Date)
	}
}