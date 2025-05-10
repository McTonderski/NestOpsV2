package middleware

import (
	"net/http"
	"strings"
	"time"

	"homelab-monitor/auth"
	"homelab-monitor/database"
	"homelab-monitor/models"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid auth header"})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ValidateAccessToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		// Save claims to context
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || role != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}

func TrackDailyApiCall() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		go func() {
			today := time.Now().Truncate(24 * time.Hour)
			var entry models.DailyApiCall
			result := database.DB.Where("date = ?", today).First(&entry)
			if result.Error != nil {
				entry = models.DailyApiCall{
					Date:  today,
					Count: 1,
				}
				database.DB.Create(&entry)
			} else {
				database.DB.Model(&entry).Update("count", entry.Count+1)
			}
		}()
	}
}
