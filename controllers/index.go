package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
)

func CombineRoutes(engine *gin.Engine) {
	homeController := &HomeController{}
	homeController.applyRoutes(engine)

	registerController := &RegisterController{DB: db.DB}
	registerController.applyRoutes(engine)
}
