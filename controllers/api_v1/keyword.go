package api_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type KeywordAPIController struct{}

func (kapi *KeywordAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/keyword", kapi.uploadKeyword)
}

func (kapi *KeywordAPIController) uploadKeyword(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
