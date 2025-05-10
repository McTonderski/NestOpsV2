package routes

import (
	"encoding/json"
	"fmt"
	"homelab-monitor/database"
	"homelab-monitor/middleware"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

var runningServicePorts = make(map[string]string)

// @Summary		List available services
// @Description	Scans the services directory and lists folders containing a docker-compose.yml
// @Tags			services
// @Produce		json
// @Success		200	{object}	ServiceList
// @Router			/services [get]
func RegisterServiceRoutes(r *gin.RouterGroup) {
	r.GET("/services", ListServices)
	r.POST("/services/:name", ToggleService)
	r.GET("/metrics/api-calls", GetDailyApiCalls)
	r.POST("/services/add", middleware.RequireAdmin(), AddService)
}

type ServiceInfo struct {
	Name   string `json:"name"`
	Status string `json:"status"`
	Port   string `json:"port,omitempty"`
}

type ServiceList struct {
	Services []ServiceInfo `json:"services"`
}

func ListServices(c *gin.Context) {
	baseDir := "./services"
	var result []ServiceInfo

	entries, err := os.ReadDir(baseDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read services directory"})
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		serviceName := entry.Name()
		servicePath := filepath.Join(baseDir, serviceName)

		composePaths := []string{
			filepath.Join(servicePath, "docker-compose.yml"),
			filepath.Join(servicePath, "docker-compose.yaml"),
			filepath.Join(servicePath, "docker", "docker-compose.yml"),
			filepath.Join(servicePath, "docker", "docker-compose.yaml"),
		}

		foundCompose := false
		for _, path := range composePaths {
			if fileExists(path) {
				foundCompose = true
				break
			}
		}

		if !foundCompose {
			continue
		}

		status := detectServiceStatus(serviceName)
		if status == "running" {
			port, ok := runningServicePorts[serviceName]
			if !ok {
				port = detectServicePort(serviceName)
				runningServicePorts[serviceName] = port
			}
			result = append(result, ServiceInfo{Name: serviceName, Status: status, Port: port})
		} else {
			result = append(result, ServiceInfo{Name: serviceName, Status: status})
		}
	}

	c.JSON(http.StatusOK, ServiceList{Services: result})
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func detectServiceStatus(serviceName string) string {
	cmd := exec.Command("docker", "compose", "ls", "--format", "json")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	var services []struct {
		Name   string `json:"Name"`
		Status string `json:"Status"`
	}
	if err := json.Unmarshal(output, &services); err != nil {
		return "unknown"
	}

	for _, svc := range services {
		if svc.Name == serviceName {
			switch {
			case strings.Contains(svc.Status, "running"):
				return "running"
			case strings.Contains(svc.Status, "exited") || strings.Contains(svc.Status, "stopped"):
				return "stopped"
			default:
				return "partial"
			}
		}
	}
	return "off"
}

func findComposeFile(servicePath string) (string, error) {
	composePaths := []string{
		filepath.Join(servicePath, "docker-compose.yml"),
		filepath.Join(servicePath, "docker-compose.yaml"),
		filepath.Join(servicePath, "docker", "docker-compose.yml"),
		filepath.Join(servicePath, "docker", "docker-compose.yaml"),
	}
	for _, path := range composePaths {
		if fileExists(path) {
			return path, nil
		}
	}
	return "", fmt.Errorf("docker-compose file not found")
}

func StartService(c *gin.Context) {
	serviceName := c.Param("name")
	baseDir := "./services"
	servicePath := filepath.Join(baseDir, serviceName)

	composeFile, err := findComposeFile(servicePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cmd := exec.Command("docker-compose", "-p", serviceName, "-f", composeFile, "up", "-d")

	output, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(output)
		if strings.Contains(msg, "no matching manifest for linux/arm64/v8") {
			// Retry with platform override
			cmd = exec.Command("docker-compose", "-p", serviceName, "-f", composeFile, "up", "-d")
			cmd.Env = append(os.Environ(), "DOCKER_DEFAULT_PLATFORM=linux/amd64")
			output, err = cmd.CombinedOutput()
			if err != nil {
				msg = string(output)
				msg += "\nRetried with DOCKER_DEFAULT_PLATFORM=linux/amd64 but still failed."
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "failed to start service",
					"details": msg,
				})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error":   "failed to start service",
				"details": msg,
			})
			return
		}
	}

	port := detectServicePort(serviceName)
	runningServicePorts[serviceName] = port

	c.JSON(http.StatusOK, gin.H{
		"message": "service started",
		"output":  string(output),
	})
}

func ToggleService(c *gin.Context) {
	serviceName := c.Param("name")
	status := detectServiceStatus(serviceName)

	if status == "stopped" || status == "off" || status == "unknown" {
		StartService(c)
	} else if status == "running" || status == "partial" {
		StopService(c)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to determine service status",
		})
	}
}

func StopService(c *gin.Context) {
	serviceName := c.Param("name")
	baseDir := "./services"
	servicePath := filepath.Join(baseDir, serviceName)

	composeFile, err := findComposeFile(servicePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	cmd := exec.Command("docker-compose", "-f", composeFile, "down")

	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "failed to stop service",
			"details": string(output),
		})
		return
	}

	delete(runningServicePorts, serviceName)

	c.JSON(http.StatusOK, gin.H{
		"message": "service stopped",
		"output":  string(output),
	})
}

func GetDailyApiCalls(c *gin.Context) {
	results, err := database.GetDailyApiCalls()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve API calls"})
		return
	}

	c.JSON(http.StatusOK, results)
}

func AddService(c *gin.Context) {
	type Request struct {
		URL string `json:"url"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil || req.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	baseDir := "./services"
	parts := strings.Split(strings.TrimSuffix(req.URL, ".git"), "/")
	serviceName := parts[len(parts)-1]
	servicePath := filepath.Join(baseDir, serviceName)

	if _, err := os.Stat(servicePath); !os.IsNotExist(err) {
		c.JSON(http.StatusConflict, gin.H{"error": "Service already exists"})
		return
	}

	cmd := exec.Command("git", "clone", req.URL, servicePath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clone repository",
			"details": string(output),
		})
		return
	}
	// Validate docker-compose file exists
	composePaths := []string{
		filepath.Join(servicePath, "docker-compose.yml"),
		filepath.Join(servicePath, "docker-compose.yaml"),
		filepath.Join(servicePath, "docker", "docker-compose.yml"),
		filepath.Join(servicePath, "docker", "docker-compose.yaml"),
	}

	foundCompose := false
	for _, path := range composePaths {
		if fileExists(path) {
			foundCompose = true
			break
		}
	}

	if !foundCompose {
		os.RemoveAll(servicePath)
		c.JSON(http.StatusBadRequest, gin.H{"error": "No docker-compose file found in the repository. Service has been removed."})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Service added successfully",
		"name":    serviceName,
	})
}

func detectServicePort(serviceName string) string {
	cmd := exec.Command("docker", "inspect", serviceName)
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	var containers []struct {
		NetworkSettings struct {
			Ports map[string][]struct {
				HostPort string `json:"HostPort"`
			} `json:"Ports"`
		} `json:"NetworkSettings"`
	}

	if err := json.Unmarshal(output, &containers); err != nil || len(containers) == 0 {
		return ""
	}

	ports := containers[0].NetworkSettings.Ports
	for _, bindings := range ports {
		for _, binding := range bindings {
			port := binding.HostPort
			resp, err := http.Get("http://localhost:" + port)
			if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 500 {
				return port
			}
		}
	}

	return ""
}
