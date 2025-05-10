package database

import (
	"homelab-monitor/models"
	"log"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("homelab.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database")
	}
	// Ensure the schema is up to date with all fields in the User model
	if err := DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("failed to migrate User schema: %v", err)
	}
	if err := DB.AutoMigrate(&models.UserAuth{}); err != nil {
		log.Fatalf("failed to migrate UserAuth schema: %v", err)
	}
	if err := DB.AutoMigrate(&models.SpeedTest{}); err != nil {
		log.Fatalf("failed to migrate SpeedTest schema: %v", err)
	}
	if err := DB.AutoMigrate(&models.DailyApiCall{}); err != nil {
		log.Fatalf("failed to migrate DailyApiCall schema: %v", err)
	}
}

func GetAllUsers() ([]models.User, error) {
	var users []models.User
	result := DB.Find(&users)
	return users, result.Error
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func UpdateUserPhoto(email, photoPath string) error {
	return DB.Model(&models.User{}).Where("email = ?", email).Update("photo_url", photoPath).Error
}

func UpdateUserDetails(email, birthday, city, street, postalCode, taxID, phone, first, last, xLink, fbLink, gitHubLink string) error {
	return DB.Model(&models.User{}).Where("email = ?", email).Updates(models.User{
		Birthday:   birthday,
		City:       city,
		Street:     street,
		PostalCode: postalCode,
		TaxID:      taxID,
		Phone:      phone,
		First:      first,
		Last:       last,
		XLink:      xLink,
		FBLink:     fbLink,
		GitHubLink: gitHubLink,
	}).Error
}

func GetUserPhotoPath(email string) (string, error) {
	var user models.User
	result := DB.First(&user, "email = ?", email)
	if result.Error != nil {
		return "", result.Error
	}
	return user.PhotoURL, nil
}

// GetDailyApiCalls retrieves the past 14 days of DailyApiCall records.
func GetDailyApiCalls() ([]models.DailyApiCall, error) {
	var results []models.DailyApiCall
	cutoff := time.Now().AddDate(0, 0, -13).Truncate(24 * time.Hour)
	result := DB.Where("date >= ?", cutoff).Order("date ASC").Find(&results)
	return results, result.Error
}
