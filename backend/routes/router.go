package routes

import (
	"homelab-monitor/controllers"
	"homelab-monitor/middleware"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // 👈 must be exact origin
		AllowMethods:     []string{"POST", "GET", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true, // 👈 must be true to send cookies or headers
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.TrackDailyApiCall())

	// --- Group API Routes under /api ---
	api := r.Group("/api")
	{
		if os.Getenv("GIN_MODE") != "release" {
			api.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		}
		api.GET("/", controllers.Dashboard)
		RegisterAuthRoutes(api)
		RegisterServiceRoutes(api)
		RegisterSystemRoutes(api)
		RegisterUserRoutes(api)
	}

	return r
}
