package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

const (
	logoutSuccessFlash = "You've been logged out"
)

type LogoutController struct{}

func (l *LogoutController) applyRoutes(engine *gin.Engine) {
	engine.POST("/logout", l.logout)
}

func (l *LogoutController) logout(c *gin.Context) {
	session.Delete(c, "user_id")

	session.AddFlash(c, logoutSuccessFlash, html.FlashNoticeKey)
	c.Redirect(http.StatusFound, "/")
}
