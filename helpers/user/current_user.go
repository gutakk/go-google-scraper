package helpers

import (
	"strconv"

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

func GetCurrentUserID(c *gin.Context) uint {
	currentUserID, exist := c.Get("currentUserID")

	if exist {
		v, ok := currentUserID.(string)
		if ok {
			userID, err := strconv.Atoi(v)
			if err != nil {
				return 0
			}
			return uint(userID)
		}
	}

	return 0
}
