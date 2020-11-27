package helpers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func RenderWithError(c *gin.Context, status int, view string, title string, err error, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title":  title,
		"errors": err.Error(),
		"data":   data,
	})
}
