package helpers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/helpers/str"
)

func RenderWithError(c *gin.Context, status int, view string, title string, err error, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title": title,
		"error": str.CapitalizeFirst(err.Error()),
		"data":  data,
	})
}

func RenderWithFlash(c *gin.Context, status int, view string, title string, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title":         title,
		"noticeFlashes": session.Flashes(c, "notice"),
		"errorFlashes":  session.Flashes(c, "error"),
		"data":          data,
	})
}
