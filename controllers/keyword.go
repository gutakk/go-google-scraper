package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	keywordTitle = "Keyword"
	keywordView  = "keyword"

	uploadSuccessFlash = "CSV uploaded successfully"
)

type KeywordController struct{}

type UploadFileForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (k *KeywordController) applyRoutes(engine *gin.RouterGroup) {
	engine.GET("/keyword", k.displayKeyword)
	engine.POST("/keyword", k.uploadKeyword)
}

func (k *KeywordController) displayKeyword(c *gin.Context) {
	html.RenderWithFlash(c, http.StatusOK, keywordView, keywordTitle, nil)
}

func (k *KeywordController) uploadKeyword(c *gin.Context) {
	form := &UploadFileForm{}
	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, errorHandler.ValidationErrorMessage(fieldErr), nil)
			return
		}
	}

	// Validate if file is CSV type
	if err := models.ValidateFileType(form.File.Header["Content-Type"][0]); err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
		return
	}

	record := models.UploadFile(c, form.File)

	// Validate if CSV has row between 1 and 1,000
	if err := models.ValidateCSVLength(len(record)); err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
		return
	}

	userID := session.Get(c, "user_id")

	var user models.User
	if userID != nil {
		user, _ = models.FindUserByID(userID)
	}

	// Save keywords to database
	_, err := models.SaveKeywords(user.ID, record)
	if err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
		return
	}

	html.RenderWithNotice(c, http.StatusOK, keywordView, keywordTitle, uploadSuccessFlash, nil)
}
