package config

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/controllers"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("templates/*")
	router.Static("/dist", "./dist")
	controllers.CombineRoutes(router)

	return router
}
