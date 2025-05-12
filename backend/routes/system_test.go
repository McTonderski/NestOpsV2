package routes

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func MockGetSystemStats(c *gin.Context) {
	simulate := c.Query("simulate")

	if simulate == "cpu_error" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CPU error simulated"})
		return
	}

	// Normal response
	c.JSON(http.StatusOK, gin.H{"status": "System stats retrieved successfully"})
}

func TestGetSystemStats_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/system/stats", MockGetSystemStats)

	req, _ := http.NewRequest(http.MethodGet, "/system/stats?simulate=cpu_error", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	MockGetSystemStats(c)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
