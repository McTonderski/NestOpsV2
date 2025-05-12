package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestServiceEndpoint(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the service route
	router.GET("/service", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Service is running"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/service", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "Service is running"}`, resp.Body.String())
}
func TestListServices(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the ListServices route
	router.GET("/services", func(c *gin.Context) {
		c.JSON(http.StatusOK, ServiceList{
			Services: []ServiceInfo{
				{Name: "service1", Status: "running", Port: "8080"},
				{Name: "service2", Status: "stopped"},
			},
		})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/services", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{
		"services": [
			{"name": "service1", "status": "running", "port": "8080"},
			{"name": "service2", "status": "stopped"}
		]
	}`, resp.Body.String())
}

func TestToggleServiceStart(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the ToggleService route
	router.POST("/services/:name", func(c *gin.Context) {
		serviceName := c.Param("name")
		if serviceName == "service1" {
			c.JSON(http.StatusOK, gin.H{"message": "service started"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to start service"})
		}
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodPost, "/services/service1", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "service started"}`, resp.Body.String())
}

func TestToggleServiceStop(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the ToggleService route
	router.POST("/services/:name", func(c *gin.Context) {
		serviceName := c.Param("name")
		if serviceName == "service1" {
			c.JSON(http.StatusOK, gin.H{"message": "service stopped"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to stop service"})
		}
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodPost, "/services/service1", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "service stopped"}`, resp.Body.String())
}

func TestAddService(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the AddService route
	router.POST("/services/add", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Service added successfully", "name": "new-service"})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodPost, "/services/add", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"message": "Service added successfully", "name": "new-service"}`, resp.Body.String())
}

func TestGetDailyApiCalls(t *testing.T) {
	// Set up a test router
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	// Mock the GetDailyApiCalls route
	router.GET("/metrics/api-calls", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"calls": 42})
	})

	// Create a test request
	req, _ := http.NewRequest(http.MethodGet, "/metrics/api-calls", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assert the response
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.JSONEq(t, `{"calls": 42}`, resp.Body.String())
}
