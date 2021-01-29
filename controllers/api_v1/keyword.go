package api_v1

import (
	"errors"
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	helpers "github.com/gutakk/go-google-scraper/helpers/user"
	"github.com/gutakk/go-google-scraper/serializers"
	"github.com/gutakk/go-google-scraper/services/keyword_service"

	"github.com/gin-gonic/gin"
)

const (
	invalidFileErr = "invalid file"
)

type KeywordAPIController struct{}

func (kapi *KeywordAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.GET("/keywords/:keyword_id", kapi.fetchKeyword)
	engine.GET("/keywords", kapi.fetchKeywords)
	engine.POST("/keywords", kapi.uploadKeyword)
}

func (kapi *KeywordAPIController) fetchKeyword(c *gin.Context) {
	currentUserID := helpers.GetCurrentUserID(c)
	keywordService := keyword_service.KeywordService{CurrentUserID: currentUserID}
	keywordID := c.Param("keyword_id")

	keyword, err := keywordService.GetKeywordResult(keywordID)
	if err != nil {
		errorResponse := api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	c.JSON(http.StatusOK, keyword_api_service.JSONAPIFormatKeywordResponse(keyword))
}

func (kapi *KeywordAPIController) fetchKeywords(c *gin.Context) {
	currentUserID := helpers.GetCurrentUserID(c)
	keywordService := keyword_service.KeywordService{CurrentUserID: currentUserID, QueryString: c.Request.URL.Query()}
	conditions := keywordService.FilterValidConditions()
	keywords, err := keywordService.GetKeywords(conditions)

	if err != nil {
		errorResponse := api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	keywordsSerializer := serializers.KeywordsSerializer{Keywords: keywords}

	c.JSON(http.StatusOK, keywordsSerializer.JSONAPIFormat())
}

func (kapi *KeywordAPIController) uploadKeyword(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: errors.New(invalidFileErr).Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse.NewErrorResponse())
		return
	}

	currentUserID := helpers.GetCurrentUserID(c)
	keywordService := keyword_service.KeywordService{CurrentUserID: currentUserID}

	err = keywordService.ValidateFileType(file.Header["Content-Type"][0])
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	filename := keywordService.UploadFile(c, file)

	parsedKeywordList, err := keywordService.ReadFile(filename)
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusUnprocessableEntity,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	// Validate if CSV has row between 1 and 1,000
	err = keywordService.ValidateCSVLength(len(parsedKeywordList))
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	// Save keywords to database
	err = keywordService.Save(parsedKeywordList)
	if err != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: err.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
