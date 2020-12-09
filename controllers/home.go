package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	helpers "github.com/gutakk/go-google-scraper/helpers/user"
)

const (
	homeTitle = "Home"
	homeView  = "home"
)

type HomeController struct{}

func (h *HomeController) applyRoutes(engine *gin.Engine) {
	engine.GET("/", h.displayHome)
}

func (h *HomeController) displayHome(c *gin.Context) {
	currentUser := helpers.GetCurrentUser(c)

	data := map[string]interface{}{
		"authenticatedUser": currentUser.ID,
		"email":             currentUser.Email,
	}

	html.RenderWithFlash(c, http.StatusOK, homeView, homeTitle, data)
}
