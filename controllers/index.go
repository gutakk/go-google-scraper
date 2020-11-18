package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
)

func CombineRoutes(engine *gin.Engine) {
	authController := &AuthController{DB: db.DB}
	authController.applyRoutes(engine)

	homeController := &HomeController{}
	homeController.applyRoutes(engine)
}
