package models

import (
	"time"

	"gorm.io/gorm"
)

type DailyApiCall struct {
	gorm.Model
	Date  time.Time `gorm:"unique;type:date" json:"date"`
	Count int       `gorm:"default:0" json:"count"`
}
