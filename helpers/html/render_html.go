package helpers

import (
	"strings"

	session "github.com/gutakk/go-google-scraper/helpers/session"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

func RenderWithError(c *gin.Context, status int, view string, title string, err error, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title": title,
		"view":  strings.ReplaceAll(view, "_", "-"),
		"error": err.Error(),
		"data":  data,
	})
}

func RenderWithFlash(c *gin.Context, status int, view string, title string, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title":         title,
		"view":          strings.ReplaceAll(view, "_", "-"),
		"noticeFlashes": session.Flashes(c, "notice"),
		"errorFlashes":  session.Flashes(c, "error"),
		"data":          data,
	})
}

func RenderErrorPage(c *gin.Context, status int, view string, title string) {
	ginview.HTML(c, status, view, gin.H{
		"title": title,
		"view":  strings.ReplaceAll(view, "_", "-"),
	})
}
