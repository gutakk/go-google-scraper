package helpers

import (
	"github.com/gin-gonic/gin"
)

func HtmlWithError(c *gin.Context, title string, view string, status int, errorMsg string) {
	c.HTML(status, view, gin.H{
		"title":  title,
		"errors": errorMsg,
	})
}

func HtmlWithNotice(c *gin.Context, title string, view string, status int, noticeMsg interface{}) {
	c.HTML(status, view, gin.H{
		"title":   title,
		"notices": noticeMsg,
	})
}
