package api_v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DummyAPIController struct{}

func (d *DummyAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.GET("/dummy", d.dummyAPI)
}

func (d *DummyAPIController) dummyAPI(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"hello": "world",
	})
}
