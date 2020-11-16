package config

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Home",
		})
	})

	health := new(controllers.HealthController)
	router.GET("/health", health.Status)

	return router
}
