package controllers

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

type HomeController struct{}

func (h *HomeController) applyRoutes(engine *gin.Engine) {
	engine.GET("/", h.displayHome)
}

func (h *HomeController) displayHome(c *gin.Context) {
	userID := session.Get(c, "user_id")

	var user models.User
	if userID != nil {
		user, _ = models.FindOneUserByID(userID)
	}

	ginview.HTML(c, http.StatusOK, "home", gin.H{
		"title":             "Home",
		"authenticatedUser": userID,
		"email":             user.Email,
	})
}
