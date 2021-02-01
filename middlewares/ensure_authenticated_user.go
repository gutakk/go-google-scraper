package middlewares

import (
	"net/http"

	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"

	"github.com/gin-gonic/gin"
)

func EnsureAuthenticatedUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID == nil {
		redirectToLogin(c)
	} else {
		_, err := models.FindUserByID(userID)

		if err != nil {
			session.Delete(c, "user_id")
			redirectToLogin(c)
		}
	}
}

func redirectToLogin(c *gin.Context) {
	session.AddFlash(c, "Login required", "error")
	c.Redirect(http.StatusFound, "/login")
	c.Abort()
}
