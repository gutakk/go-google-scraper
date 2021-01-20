package html

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/helpers/session"
)

func RenderWithError(c *gin.Context, status int, view string, title string, err error, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title": title,
		"view":  view,
		"error": err.Error(),
		"data":  data,
	})
}

func RenderWithFlash(c *gin.Context, status int, view string, title string, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title":         title,
		"view":          view,
		"noticeFlashes": session.Flashes(c, "notice"),
		"errorFlashes":  session.Flashes(c, "error"),
		"data":          data,
	})
}
