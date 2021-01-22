package api_v1

import (
	"github.com/gutakk/go-google-scraper/services/oauth_service"

	"github.com/gin-gonic/gin"
)

type TokenAPIController struct{}

// TODO: Unit test in login API PR as grant type need to change to password
func (t *TokenAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/token", oauth_service.GenerateToken)
}
