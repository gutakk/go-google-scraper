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
	keywordService := initKeywordService(c)
	data := getKeywordsData(keywordService)

	html.RenderWithFlash(c, http.StatusOK, keywordView, keywordTitle, data)
}

func (k *KeywordController) uploadKeyword(c *gin.Context) {
	keywordService := initKeywordService(c)
	data := getKeywordsData(keywordService)

	form := &UploadFileForm{}
	if err := c.ShouldBind(form); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, errorHandler.ValidationErrorMessage(fieldErr), data)
			return
		}
	}

	// Validate if file is CSV type
	validateTypeErr := keywordService.ValidateFileType(form.File.Header["Content-Type"][0])
	if validateTypeErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, validateTypeErr, data)
		return
	}

	filename := keywordService.UploadFile(c, form.File)

	record, readFileErr := keywordService.ReadFile(filename)
	if readFileErr != nil {
		html.RenderWithError(c, http.StatusUnprocessableEntity, keywordView, keywordTitle, readFileErr, data)
	}

	// Validate if CSV has row between 1 and 1,000
	validateLengthErr := keywordService.ValidateCSVLength(len(record))
	if validateLengthErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, validateLengthErr, data)
		return
	}

	// Save keywords to database
	_, saveKeywordsErr := keywordService.Save(record)
	if saveKeywordsErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, saveKeywordsErr, data)
		return
	}

	session.AddFlash(c, uploadSuccessFlash, "notice")
	c.Redirect(http.StatusFound, "/keyword")
}

func getKeywordsData(keywordService keyword_service.KeywordService) map[string]interface{} {
	keywords, _ := keywordService.GetAll()

	return map[string]interface{}{
		"keywords": keywords,
	}
}

func initKeywordService(c *gin.Context) keyword_service.KeywordService {
	currentUser := helpers.GetCurrentUser(c)
	return keyword_service.KeywordService{CurrentUserID: currentUser.ID}
}
