package helpers

import (
	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
)

const (
	FlashNoticeKey = "flashNotices"
	FlashErrorKey  = "flashErrors"
)

func RenderWithError(c *gin.Context, status int, view string, title string, err error, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title":  title,
		"errors": err.Error(),
		"data":   data,
	})
}

func RenderWithFlash(c *gin.Context, status int, view string, title string, data map[string]interface{}) {
	flashKey, flashValue := getFlashMessage(c)

	ginview.HTML(c, status, view, gin.H{
		"title":  title,
		flashKey: flashValue,
		"data":   data,
	})
}

func Render(c *gin.Context, status int, view string, title string, data map[string]interface{}) {
	ginview.HTML(c, status, view, gin.H{
		"title": title,
		"data":  data,
	})
}

func getFlashMessage(c *gin.Context) (string, interface{}) {
	flashNotices := session.Flashes(c, FlashNoticeKey)
	if flashNotices != nil {
		return FlashNoticeKey, flashNotices
	}

	flashErrors := session.Flashes(c, FlashErrorKey)
	return FlashErrorKey, flashErrors
}
