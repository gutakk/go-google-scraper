package middlewares

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
)

func ValidateToken(c *gin.Context) {
	_, err := oauth.GetOAuthServer().ValidationBearerToken(c.Request)
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Status: http.StatusUnauthorized,
		}

		if err.Error() == "sql: no rows in result set" {
			errorResponse.Detail = errors.ErrInvalidAccessToken.Error()
		} else {
			errorResponse.Detail = err.Error()
		}

		c.JSON(http.StatusUnauthorized, errorResponse.ConstructErrorResponse())
		c.Abort()
	}
}
