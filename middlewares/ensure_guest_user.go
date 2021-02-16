package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func EnsureGuestUser(c *gin.Context) {
	userID := c.MustGet("currentUser")

	if userID != nil {
		c.Redirect(http.StatusFound, "/")
		c.Abort()
	}
}
