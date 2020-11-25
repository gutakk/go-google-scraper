package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
)

func CombineRoutes(engine *gin.Engine) {
	homeController := &HomeController{}
	homeController.applyRoutes(engine)

	loginController := &LoginController{DB: db.DB}
	loginController.applyRoutes(engine)

	registerController := &RegisterController{}
	registerController.applyRoutes(engine)
}
