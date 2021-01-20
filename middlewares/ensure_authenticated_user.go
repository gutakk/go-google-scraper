package middlewares

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/session"

	"github.com/gin-gonic/gin"
)

func EnsureAuthenticatedUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID == nil {
		session.AddFlash(c, "Login required", "error")
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
	}
}
