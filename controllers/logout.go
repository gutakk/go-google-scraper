package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

const (
	logoutSuccessFlash = "You have been logged out"
)

type LogoutController struct{}

func (l *LogoutController) applyRoutes(engine *gin.RouterGroup) {
	engine.POST("/logout", l.logout)
}

func (l *LogoutController) logout(c *gin.Context) {
	session.Delete(c, "user_id")

	session.AddFlash(c, logoutSuccessFlash, "notice")
	c.Redirect(http.StatusFound, "/")
}
