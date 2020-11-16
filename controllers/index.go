package controllers

import "github.com/gin-gonic/gin"

func CombineRoutes(engine *gin.Engine) {
	new(AuthController).applyRoutes(engine)
	new(HomeController).applyRoutes(engine)
}
