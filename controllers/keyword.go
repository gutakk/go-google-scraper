package controllers

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"path/filepath"

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
	filename := "dist/" + filepath.Base(file.Filename)
	_ = c.SaveUploadedFile(file, filename)

	csvfile, _ := os.Open(filename)

	r := csv.NewReader(csvfile)
	record, _ := r.ReadAll()
	log.Printf("!!!!!!!!!!!!!!!! %v", len(record))

	for i, v := range record {
		log.Printf("################## %v", i)
		log.Printf("################## %v", v[0])
	}
}
