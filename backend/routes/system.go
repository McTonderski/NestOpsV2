package routes

import (
	"errors"
	"homelab-monitor/database"
	"homelab-monitor/middleware"
	"homelab-monitor/models"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

type SystemStats struct {
	CPUUsagePercent float64                `json:"cpu_usage_percent"`
	RAM             *mem.VirtualMemoryStat `json:"ram"`
	Disk            *disk.UsageStat        `json:"disk"`
	Temperature     float64                `json:"temperature,omitempty"`
}

// @Summary	Get system statistics
// @Tags		system
// @Produce	json
// @Success	200	{object}	SystemStats
// @Router		/api/system/stats [get]
func RegisterSystemRoutes(r *gin.RouterGroup) {
	r.GET("/system/stats", middleware.RequireAuth(), GetSystemStats)
	r.GET("/system/speedtest", middleware.RequireAuth(), GetSpeedTestData)
}

func readTemperatureFromSys() (float64, error) {
	var maxTemp float64
	files, err := filepath.Glob("/sys/class/thermal/thermal_zone*/temp")
	if err != nil {
		return 0, err
	}

	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		temp, err := strconv.ParseFloat(string(data), 64)
		if err != nil {
			continue
		}
		if temp > maxTemp {
			maxTemp = temp
		}
	}
	if maxTemp == 0 {
		return 0, errors.ErrUnsupported
	}
	return maxTemp / 1000, nil // Convert from millidegree Celsius to degree Celsius
}

func GetSystemStats(c *gin.Context) {
	// CPU
	cpuPercent, err := cpu.Percent(0, false)
	if err != nil || len(cpuPercent) == 0 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get CPU usage"})
		return
	}

	// RAM
	vmStats, err := mem.VirtualMemory()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get memory stats"})
		return
	}

	// Disk
	diskStats, err := disk.Usage("/")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get disk stats"})
		return
	}

	// Temp (using readTemperatureFromSys)
	temperature, err := readTemperatureFromSys()
	if err != nil {
		log.Printf("Could not read temperature: %v\n", err)
		temperature = -1
	}

	stats := SystemStats{
		CPUUsagePercent: cpuPercent[0],
		RAM:             vmStats,
		Disk:            diskStats,
		Temperature:     temperature,
	}
	c.JSON(http.StatusOK, stats)
}

func GetSpeedTestData(c *gin.Context) {
	var results []models.SpeedTest
	if err := database.DB.Order("created_at desc").Limit(250).Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve speed test data"})
		return
	}
	c.JSON(http.StatusOK, results)
}
