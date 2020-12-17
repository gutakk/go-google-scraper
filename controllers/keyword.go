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
	"github.com/gutakk/go-google-scraper/presenters"
	"github.com/gutakk/go-google-scraper/services/keyword_service"
)

const (
	keywordTitle      = "Keyword"
	keywordView       = "keyword"
	keywordResultView = "keyword_result"

	uploadSuccessFlash = "CSV uploaded successfully"
)

type KeywordController struct{}

type UploadFileForm struct {
	File *multipart.FileHeader `form:"file" binding:"required"`
}

func (k *KeywordController) applyRoutes(engine *gin.RouterGroup) {
	engine.GET("/keyword/:keyword_id", k.displayKeywordResult)
	engine.GET("/keyword", k.displayKeyword)
	engine.POST("/keyword", k.uploadKeyword)
}

func (k *KeywordController) displayKeywordResult(c *gin.Context) {
	keywordService := initKeywordService(c)
	keywordID := c.Param("keyword_id")
	data, err := getKeywordData(keywordService, keywordID)

	if err != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordResultView, keywordTitle, err, data)
		return
	}

	html.RenderWithFlash(c, http.StatusOK, keywordResultView, keywordTitle, data)
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

	parsedKeywordList, readFileErr := keywordService.ReadFile(filename)
	if readFileErr != nil {
		html.RenderWithError(c, http.StatusUnprocessableEntity, keywordView, keywordTitle, readFileErr, data)
	}

	// Validate if CSV has row between 1 and 1,000
	validateLengthErr := keywordService.ValidateCSVLength(len(parsedKeywordList))
	if validateLengthErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, validateLengthErr, data)
		return
	}

	// Save keywords to database
	saveKeywordsErr := keywordService.Save(parsedKeywordList)
	if saveKeywordsErr != nil {
		html.RenderWithError(c, http.StatusBadRequest, keywordView, keywordTitle, saveKeywordsErr, data)
		return
	}

	session.AddFlash(c, uploadSuccessFlash, "notice")
	c.Redirect(http.StatusFound, "/keyword")
}

func getKeywordData(keywordService keyword_service.KeywordService, keywordID string) (map[string]interface{}, error) {
	keyword, err := keywordService.GetSingle(keywordID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"keyword": keyword,
	}, nil
}

func getKeywordsData(keywordService keyword_service.KeywordService) map[string]interface{} {
	keywords, _ := keywordService.GetAll()
	var keywordPresenters []presenters.KeywordPresenter

	for _, k := range keywords {
		keywordPresenters = append(keywordPresenters, presenters.KeywordPresenter{Keyword: k})
	}

	return map[string]interface{}{
		"keywordPresenters": keywordPresenters,
	}
}

func initKeywordService(c *gin.Context) keyword_service.KeywordService {
	currentUser := helpers.GetCurrentUser(c)
	return keyword_service.KeywordService{CurrentUserID: currentUser.ID}
}
