package middlewares

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/session"

	"github.com/gin-gonic/gin"
)

func EnsureGuestUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}
}
