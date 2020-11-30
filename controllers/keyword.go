package controllers

import (
	"encoding/csv"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	"github.com/gutakk/go-google-scraper/models"
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

	if err := models.ValidateFileFormat(file.Header["Content-Type"][0]); err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
	}

	filename := "dist/" + filepath.Base(file.Filename)
	_ = c.SaveUploadedFile(file, filename)

	csvfile, _ := os.Open(filename)

	r := csv.NewReader(csvfile)
	record, _ := r.ReadAll()
	if err := models.ValidateCSVLength(len(record)); err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
	}

	for i, v := range record {
		log.Printf("################## %v", i)
		log.Printf("################## %v", v[0])
	}
}