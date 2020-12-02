package tests

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/middlewares"
)

// Helper function to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	router := gin.Default()
	router = middlewares.SetupMiddlewares(router)

	if withTemplates {
		router.HTMLRender = ginview.New(config.AppGoviewConfig())
		router.Static("/dist", "./dist")
	}

	return router
}
