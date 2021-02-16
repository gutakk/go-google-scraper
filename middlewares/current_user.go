package middlewares

import (
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	userID := session.Get(c, "user_id")

	if userID == nil {
		c.Set("currentUser", nil)
		return
	}

	var user models.User
	var err error
	user, err = models.FindUserByID(userID)

	if err != nil {
		c.Set("currentUser", nil)
		return
	}

	c.Set("currentUser", user)
}
