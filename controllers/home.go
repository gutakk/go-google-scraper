package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

type HomeController struct{}

func (h *HomeController) applyRoutes(engine *gin.Engine) {
	engine.GET("/", h.displayHome)
}

func (h *HomeController) displayHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":   "Home",
		"flashes": session.Flashes(c),
	})
}
