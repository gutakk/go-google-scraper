package middlewares

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
)

func ValidateToken(c *gin.Context) {
	_, err := oauth.GetOAuthServer().ValidationBearerToken(c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}
}
