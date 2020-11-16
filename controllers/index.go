package controllers

import "github.com/gin-gonic/gin"

func CombineRoutes(e *gin.Engine) {
	new(AuthController).applyRoutes(e)
	new(HomeController).applyRoutes(e)
}
