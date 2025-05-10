package models

import (
	"gorm.io/gorm"
	"time"
)

type SpeedTest struct {
	gorm.Model
	DownloadMbps float64   `json:"download_mbps"`
	UploadMbps   float64   `json:"upload_mbps"`
	PingMs       float64   `json:"ping_ms"`
	Timestamp    time.Time `json:"timestamp"`
}

func PruneOldSpeedTests(db *gorm.DB, maxRecords int) error {
	var ids []uint
	if err := db.Model(&SpeedTest{}).
		Select("id").
		Order("timestamp desc").
		Offset(maxRecords).
		Pluck("id", &ids).Error; err != nil {
		return err
	}
	if len(ids) > 0 {
		return db.Delete(&SpeedTest{}, ids).Error
	}
	return nil
}
