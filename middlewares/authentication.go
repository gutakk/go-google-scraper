package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

func AuthenticatedUserNotAllowed(c *gin.Context) {
	user := session.Get(c, "user_id")

	if user != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.Next()
}
