package middlewares

import (
	"net/http"

	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"

	"github.com/gin-gonic/gin"
)

func EnsureGuestUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID != nil {
		_, err := models.FindUserByID(userID)

		if err != nil {
			session.Delete(c, "user_id")
			c.Redirect(http.StatusFound, c.FullPath())
		} else {
			c.Redirect(http.StatusFound, "/")
		}

		c.Abort()
	}
}
