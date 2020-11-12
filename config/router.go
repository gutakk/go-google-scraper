package config

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	health := new(controllers.HealthController)

	router.GET("/health", health.Status)

	return router
}
