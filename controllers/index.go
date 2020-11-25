package controllers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	new(HomeController).applyRoutes(engine)

	new(RegisterController).applyRoutes(EnsureNoAuthenticationGroup(engine))
	new(UserSessionController).applyRoutes(EnsureNoAuthenticationGroup(engine))
}

func EnsureNoAuthenticationGroup(engine *gin.Engine) *gin.RouterGroup {
	mw := ginview.NewMiddleware(config.AuthenticationGoviewConfig())

	ensureNoAuthenticationGroup := engine.Group("", mw)
	ensureNoAuthenticationGroup.Use(middlewares.EnsureNoAuthentiction)

	return ensureNoAuthenticationGroup
}
