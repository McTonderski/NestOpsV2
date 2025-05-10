// main.go
package main

import (
	"fmt"
	"homelab-monitor/config"
	"homelab-monitor/database"
	"homelab-monitor/routes"
	speedtest "homelab-monitor/utils"
	"log"
	"time"
)

func main() {
	config.LoadEnv()
	database.InitDB()
	go startBackgroundTask() // periodically tasks
	r := routes.SetupRouter()
	log.Println("Server running on http://localhost:8080")
	r.Run(":8080")
}

func startBackgroundTask() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		fmt.Println("Running periodic task at", time.Now())
		_, err := speedtest.MeasureSpeed()
		if err != nil {
			log.Printf("Speedtest failed: %v", err)
		}
		<-ticker.C
	}

}
