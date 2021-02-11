package tests

import (
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"

	"github.com/gin-gonic/gin"
)

// Helper function to create a router during testing
func SetupTestRouter() *gin.Engine {
	engine := config.SetupRouter()
	controllers.CombineRoutes(engine)
	return engine
}
