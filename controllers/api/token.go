package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
)

type TokenAPIController struct{}

func (t *TokenAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/token", t.generateToken)
}

func (t *TokenAPIController) generateToken(c *gin.Context) {
	err := config.GetOAuthServer().HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(400, nil)
	}
}
