package middlewares

import (
	"net/http"

	session "github.com/gutakk/go-google-scraper/helpers/session"

	"github.com/gin-gonic/gin"
)

func EnsureAuthenticatedUser(c *gin.Context) {
	userID := c.MustGet("currentUser")

	if userID == nil {
		session.AddFlash(c, "Login required", "error")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
