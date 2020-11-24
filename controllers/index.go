package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	new(HomeController).applyRoutes(engine)

	new(RegisterController).applyRoutes(AuthenticatedUserNotAllowedGroup(engine))
	new(UserSessionController).applyRoutes(AuthenticatedUserNotAllowedGroup(engine))
}

func AuthenticatedUserNotAllowedGroup(engine *gin.Engine) *gin.RouterGroup {
	authenticatedUserNotAllowedGroup := engine.Group("")
	authenticatedUserNotAllowedGroup.Use(middlewares.AuthenticatedUserNotAllowed)

	return authenticatedUserNotAllowedGroup
}
