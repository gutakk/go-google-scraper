package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	new(HomeController).applyRoutes(engine)

	authenticatedUserNotAllowedGroup := engine.Group("")
	authenticatedUserNotAllowedGroup.Use(middlewares.AuthenticatedUserNotAllowed)

	new(RegisterController).applyRoutes(authenticatedUserNotAllowedGroup)
	new(UserSessionController).applyRoutes(authenticatedUserNotAllowedGroup)
}
