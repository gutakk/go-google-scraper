package config

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("go-google-scraper", store))

	router.HTMLRender = ginview.New(AppGoviewConfig())
	router.Static("/dist", "./dist")

	return router
}
