package controllers

import (
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	session "github.com/gutakk/go-google-scraper/helpers/session"
	helpers "github.com/gutakk/go-google-scraper/helpers/user"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/keyword_service"
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
	currentUser := helpers.GetCurrentUser(c)
	keywordService := keyword_service.Keyword{CurrentUserID: currentUser.ID}

	keywords, err := keywordService.GetAll()
	if err != nil {
		html.RenderWithError(c, http.StatusUnprocessableEntity, keywordView, keywordTitle, err, nil)
	}

	data := map[string]interface{}{
		"keywords": keywords,
	}

	html.RenderWithFlash(c, http.StatusOK, keywordView, keywordTitle, data)
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

	filename := models.UploadFile(c, form.File)
	record, readFileErr := models.ReadFile(filename)
	if readFileErr != nil {
		html.RenderWithError(c, http.StatusUnprocessableEntity, keywordView, keywordTitle, readFileErr, nil)
	}

	// Validate if CSV has row between 1 and 1,000
	if err := models.ValidateCSVLength(len(record)); err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, err, nil)
		return
	}

	currentUser := helpers.GetCurrentUser(c)
	// Save keywords to database
	_, saveKeywordsErr := models.SaveKeywords(currentUser.ID, record)
	if saveKeywordsErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, saveKeywordsErr, nil)
		return
	}

	session.AddFlash(c, uploadSuccessFlash, "notice")
	c.Redirect(http.StatusFound, "/keyword")
}
