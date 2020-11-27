package controllers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	new(HomeController).applyRoutes(engine)
	new(LogoutController).applyRoutes(engine)

	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))
}

func EnsureGuestUserGroup(engine *gin.Engine) *gin.RouterGroup {
	mw := ginview.NewMiddleware(config.AuthenticationGoviewConfig())

	ensureGuestUserGroup := engine.Group("", mw)
	ensureGuestUserGroup.Use(middlewares.EnsureGuestUser)

	return ensureGuestUserGroup
}
