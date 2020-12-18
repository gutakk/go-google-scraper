package helpers

import (
	"io/ioutil"
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	session "github.com/gutakk/go-google-scraper/helpers/session"
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

func Render404(c *gin.Context) {
	html, _ := ioutil.ReadFile("templates/404.html")
	c.Writer.WriteHeader(http.StatusNotFound)
	_, _ = c.Writer.Write(html)
}
