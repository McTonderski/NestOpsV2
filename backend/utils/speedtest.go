package speedtest

import (
	"log"
	"time"

	"homelab-monitor/database"
	"homelab-monitor/models"

	"github.com/showwin/speedtest-go/speedtest"
)

type Result struct {
	DownloadMbps float64   `json:"download_mbps"`
	UploadMbps   float64   `json:"upload_mbps"`
	LatencyMs    float64   `json:"latency_ms"`
	Timestamp    time.Time `json:"timestamp"`
}

func MeasureSpeed() (*Result, error) {
	var speedtestClient = speedtest.New()
	serverList, _ := speedtestClient.FetchServers()
	targets, _ := serverList.FindServer([]int{})
	log.Printf("Running speed test")
	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()

		result := models.SpeedTest{
			DownloadMbps: float64(s.DLSpeed) / (1024 * 1024),
			UploadMbps:   float64(s.ULSpeed) / (1024 * 1024),
			PingMs:       s.Latency.Seconds() * 1000,
			Timestamp:    time.Now(),
		}
		if res := database.DB.Create(&result); res.Error != nil {
			log.Printf("Failed to save speed test result: %v", res.Error)
			return nil, res.Error
		}

		log.Printf("Speed test recorded: Download=%.2f Mbps, Upload=%.2f Mbps, Ping=%.2f ms", result.DownloadMbps, result.UploadMbps, result.PingMs)
		models.PruneOldSpeedTests(database.DB, 250)

		s.Context.Reset()
	}

	return nil, nil
}
