package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeController struct{}

func (h *HomeController) applyRoutes(engine *gin.Engine) {
	engine.GET("/", h.displayHome)
}

func (h *HomeController) displayHome(c *gin.Context) {
	c.HTML(http.StatusOK, "home", gin.H{
		"title": "Home",
	})
}
