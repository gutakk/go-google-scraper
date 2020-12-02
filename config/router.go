package config

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	router.HTMLRender = ginview.New(AppGoviewConfig())
	router.Static("/dist", "./dist")

	return router
}
