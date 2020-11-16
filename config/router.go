package config

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.LoadHTMLGlob("templates/*")
	router.Static("/dist", "./dist")
	controllers.CombineRoutes(router)

	return router
}
