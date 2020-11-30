package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	html "github.com/gutakk/go-google-scraper/helpers/html"
)

const (
	keywordTitle = "Keyword"
	keywordView  = "keyword"
)

type KeywordController struct{}

func (k *KeywordController) applyRoutes(engine *gin.Engine) {
	engine.GET("/keyword", k.displayKeyword)
	engine.POST("/keyword", k.uploadKeyword)
}

func (k *KeywordController) displayKeyword(c *gin.Context) {
	html.RenderWithFlash(c, http.StatusOK, keywordView, keywordTitle, nil)
}

func (k *KeywordController) uploadKeyword(c *gin.Context) {
	file, _ := c.FormFile("file")
	log.Println(file.Filename)
}
