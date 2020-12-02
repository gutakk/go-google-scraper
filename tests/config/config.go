package tests

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/middlewares"
)

// Helper function to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("go-google-scraper", store))
	router.Use(middlewares.CurrentUser)

	if withTemplates {
		router.HTMLRender = ginview.New(config.AppGoviewConfig())
		router.Static("/dist", "./dist")
	}

	return router
}
