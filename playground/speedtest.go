package main

import (
	"log"
	"time"

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

	for _, s := range targets {
		s.PingTest(nil)
		s.DownloadTest()
		s.UploadTest()

		result := Result{
			DownloadMbps: float64(s.DLSpeed) / (1024 * 1024), // convert from Bps to Mbps
			UploadMbps:   float64(s.ULSpeed) / (1024 * 1024),
			LatencyMs:    s.Latency.Seconds() * 1000,
			Timestamp:    time.Now(),
		}
		// database.DB.Create(&result)
		log.Printf("Speed test recorded: Download=%.2f Mbps, Upload=%.2f Mbps, Ping=%.2f ms", result.DownloadMbps, result.UploadMbps, result.LatencyMs)
		// models.PruneOldSpeedTests(database.DB, 250)

		s.Context.Reset()
	}
	return nil, nil
}

func main() {
	MeasureSpeed()
}
