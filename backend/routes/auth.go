package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"homelab-monitor/auth"
	"homelab-monitor/database"
	"homelab-monitor/models"
)

// | Type | Lifespan | Stored Where? | Purpose |
// |------|----------|---------------|---------|
// | Access | Short (15m) | Sent via Bearer header | Auth for protected routes |
// | Refresh | Long (7d) | Stored in DB or memory | To issue new access tokens |

func RegisterAuthRoutes(r *gin.RouterGroup) {
	r.POST("/auth/refresh", RefreshToken)
	r.POST("/auth/register", Register)
	r.POST("/auth/login", Login)
}

func Register(c *gin.Context) {
	var req models.UserAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := models.UserAuth{
		Email:    req.Email,
		Password: string(hashedPassword),
		First:    string(req.First),
		Last:     string(req.Last),
		Role:     "user", // default role
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "username already exists or DB error"})
		return
	}
	err = CreateUserFromAuth(database.DB, user)
	if err != nil {
		// TODO: fix later if create user fails return error and remove entry from DB create userAuth
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

func Login(c *gin.Context) {
	var req models.UserAuth
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid json"})
		return
	}

	var user models.UserAuth
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	accessToken, err := auth.GenerateAccessToken(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate access token"})
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate refresh token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "login successful",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	var request struct {
		RefreshToken string `json:"refresh_token"`
	}

	if err := c.ShouldBindJSON(&request); err != nil || request.RefreshToken == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing refresh_token"})
		return
	}

	username, err := auth.ValidateRefreshToken(request.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	// get user's role from DB
	var user models.UserAuth
	if err := database.DB.First(&user, "email = ?", username).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	newAccessToken, err := auth.GenerateAccessToken(user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create new access token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token": newAccessToken,
	})
}

func CreateUserFromAuth(db *gorm.DB, auth models.UserAuth) error {
	user := models.User{
		Email: auth.Email,
		First: auth.First,
		Last:  auth.Last,
	}
	return db.Create(&user).Error
}
