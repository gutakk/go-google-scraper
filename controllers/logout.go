package controllers

import (
	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

type LogoutController struct{}

func (l *LogoutController) applyRoutes(engine *gin.Engine) {
	engine.POST("/logout", l.logout)
}

func (l *LogoutController) logout(c *gin.Context) {
	session.Delete(c, "user_id")
}
