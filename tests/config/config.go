package tests

import (
	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
)

// Helper function to create a router during testing
func SetupTestRouter() *gin.Engine {
	engine := config.SetupRouter()
	controllers.CombineRoutes(engine)
	return engine
}

func SetupTestOAuthServer() {
	err := oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(errorconf.StartOAuthServerFailure, err)
	}
}
