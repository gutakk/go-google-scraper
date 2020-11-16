package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct{}

func (h *HomeController) applyRoutes(e *gin.Engine) {
	e.GET("/", h.displayHome)
}

func (h *HomeController) displayHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Home",
	})
}
