package middlewares

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	pgAdapter "github.com/vgarvardt/go-pg-adapter"
)

func ValidateToken(c *gin.Context) {
	token, err := oauth.GetOAuthServer().ValidationBearerToken(c.Request)
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Status: http.StatusUnauthorized,
		}

		if err == pgAdapter.ErrNoRows {
			errorResponse.Detail = errors.ErrInvalidAccessToken.Error()
		} else {
			errorResponse.Detail = err.Error()
		}

		c.JSON(http.StatusUnauthorized, errorResponse.NewErrorResponse())
		c.Abort()
		return
	}

	c.Set("currentUserID", token.GetUserID())
}
