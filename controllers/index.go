package controllers

import (
	"github.com/gin-gonic/gin"
)

func CombineRoutes(engine *gin.Engine) {
	homeController := &HomeController{}
	homeController.applyRoutes(engine)

	userSessionController := &UserSessionController{}
	userSessionController.applyRoutes(engine)

	registerController := &RegisterController{}
	registerController.applyRoutes(engine)
}
