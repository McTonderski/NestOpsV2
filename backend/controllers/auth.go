package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	// Dummy response for now
	c.JSON(http.StatusOK, gin.H{"token": "dummy-jwt-token"})
}
