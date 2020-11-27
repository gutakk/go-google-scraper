package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
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
	userID := session.Get(c, "user_id")

	var user models.User
	if userID != nil {
		user, _ = models.FindUserByID(userID)
	}

	data := map[string]interface{}{
		"authenticatedUser": userID,
		"email":             user.Email,
	}

	html.RenderWithFlash(c, http.StatusOK, homeView, homeTitle, data)
}
