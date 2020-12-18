package controllers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	// Not found
	engine.NoRoute(func(c *gin.Context) {
		html.RenderNotFound(c)
	})

	// No group
	new(HomeController).applyRoutes(engine)

	// Ensure authenticated user group
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))
	new(LogoutController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Ensure guest user group
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))
}

func EnsureAuthenticatedUserGroup(engine *gin.Engine) *gin.RouterGroup {
	ensureAuthenticatedUserGroup := engine.Group("")
	ensureAuthenticatedUserGroup.Use(middlewares.EnsureAuthenticatedUser)

	return ensureAuthenticatedUserGroup
}

func EnsureGuestUserGroup(engine *gin.Engine) *gin.RouterGroup {
	mw := ginview.NewMiddleware(config.AuthenticationGoviewConfig())

	ensureGuestUserGroup := engine.Group("", mw)
	ensureGuestUserGroup.Use(middlewares.EnsureGuestUser)

	return ensureGuestUserGroup
}
