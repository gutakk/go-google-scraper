package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

func EnsureNoAuthentication(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Next()
}
