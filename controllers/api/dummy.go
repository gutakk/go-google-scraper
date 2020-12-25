package api

import "github.com/gin-gonic/gin"

type DummyAPIController struct{}

func (d *DummyAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.GET("/dummy", d.dummyAPI)
}

func (d *DummyAPIController) dummyAPI(c *gin.Context) {
	c.JSON(200, gin.H{
		"hello": "world",
	})
}
