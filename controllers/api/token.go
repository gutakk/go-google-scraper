package api

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/oauth"

	"github.com/gin-gonic/gin"
)

type TokenAPIController struct{}

func (t *TokenAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/token", t.generateToken)
}

func (t *TokenAPIController) generateToken(c *gin.Context) {
	err := oauth.GetOAuthServer().HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
}
