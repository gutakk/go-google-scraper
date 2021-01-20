package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/helpers/session"
)

func EnsureAuthenticatedUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID == nil {
		session.AddFlash(c, "Login required", "error")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
