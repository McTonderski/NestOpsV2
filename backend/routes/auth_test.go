package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"homelab-monitor/database"
	"homelab-monitor/models"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	api := r.Group("/api")
	RegisterAuthRoutes(api)
	return r
}

// Initialize a mock database
func setupTestDB() {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	database.DB = db

	// Migrate the schema for the UserAuth model
	database.DB.AutoMigrate(&models.UserAuth{})
}

func TestRegister(t *testing.T) {
	setupTestDB() // Initialize the mock database
	router := setupRouter()

	payload := models.UserAuth{
		Email:    "test@example.com",
		Password: "password123",
		First:    "John",
		Last:     "Doe",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "user registered successfully")
}

func TestLogin(t *testing.T) {
	setupTestDB() // Initialize the mock database
	router := setupRouter()

	// Hash the password and create a mock user in the database
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	database.DB.Create(&models.UserAuth{
		Email:    "test@example.com",
		Password: string(hashedPassword),
	})

	// Prepare the login payload
	payload := models.UserAuth{
		Email:    "test@example.com",
		Password: "password123",
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(w, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "login successful")
}
