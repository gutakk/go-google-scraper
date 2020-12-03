package helpers

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/models"
)

func GetCurrentUser(c *gin.Context) models.User {
	currentUser := c.MustGet("currentUser")

	if currentUser != nil {
		v, ok := currentUser.(models.User)
		if ok {
			return v
		}
	}
	return models.User{}
}
