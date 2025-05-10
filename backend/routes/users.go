package routes

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"homelab-monitor/database"

	"homelab-monitor/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.RouterGroup) {
	r.GET("/users", middleware.RequireAdmin(), ListUsers)
	r.GET("/users/me", middleware.RequireAuth(), GetUser)
	r.POST("/users/me/photo", middleware.RequireAuth(), UploadUserPhoto)
	r.PUT("/users/me/details", middleware.RequireAuth(), UpdateUserDetails)
	r.GET("/users/me/photo", middleware.RequireAuth(), GetUserPhoto)
}

func ListUsers(c *gin.Context) {
	users, err := database.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not retrieve users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func GetUser(c *gin.Context) {
	email := c.MustGet("username").(string)
	log.Printf("email: %s", email)
	user, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func UploadUserPhoto(c *gin.Context) {
	userID := c.MustGet("username").(string)

	// Pobierz aktualną ścieżkę do zdjęcia
	oldPhotoPath, _ := database.GetUserPhotoPath(userID)

	// Odbierz plik od użytkownika
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file"})
		return
	}

	// Zapisz nowy plik
	filePath := filepath.Join("uploads", fmt.Sprintf("%s_%s", userID, file.Filename))
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save file"})
		return
	}

	// Zaktualizuj wpis w bazie
	if err := database.UpdateUserPhoto(userID, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user photo"})
		return
	}

	// Usuń stare zdjęcie, jeśli jest inne
	if oldPhotoPath != "" && oldPhotoPath != filePath {
		_ = os.Remove(oldPhotoPath)
	}

	c.JSON(http.StatusOK, gin.H{"photoUrl": filePath})
}

func UpdateUserDetails(c *gin.Context) {
	email := c.MustGet("username").(string)

	var updates struct {
		Birthday   string `json:"birthday"`
		City       string `json:"city"`
		Street     string `json:"street"`
		PostalCode string `json:"postalCode"`
		TaxID      string `json:"taxID"`
		Phone      string `json:"phone"`
		First      string `json:"first"`
		Last       string `json:"last"`
		XLink      string `json:"xLink"`
		FBLink     string `json:"fbLink"`
		GitHubLink string `json:"gitHubLink"`
	}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := database.UpdateUserDetails(email, updates.Birthday, updates.City, updates.Street, updates.PostalCode, updates.TaxID, updates.Phone, updates.First, updates.Last, updates.XLink, updates.FBLink, updates.GitHubLink)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update user details"})
		return
	}

	updated, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func GetUserPhoto(c *gin.Context) {
	userID := c.MustGet("username").(string)

	photoPath, err := database.GetUserPhotoPath(userID)
	if err != nil || photoPath == "" {
		c.JSON(http.StatusNotFound, gin.H{"error": "Photo not found"})
		return
	}

	c.File(photoPath)
}
