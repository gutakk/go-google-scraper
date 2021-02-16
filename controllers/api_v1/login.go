package api_v1

import (
	"github.com/gutakk/go-google-scraper/services/oauth_service"

	"github.com/gin-gonic/gin"
)

type LoginAPIController struct{}

func (t *LoginAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/login", oauth_service.GenerateToken)
}
