package api_v1

import (
	"errors"
	"net/http"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	helpers "github.com/gutakk/go-google-scraper/helpers/user"
	"github.com/gutakk/go-google-scraper/services/keyword_service"

	"github.com/gin-gonic/gin"
)

const (
	invalidFileErr = "invalid file"
)

type KeywordAPIController struct{}

func (kapi *KeywordAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/keywords", kapi.uploadKeyword)
}

func (kapi *KeywordAPIController) uploadKeyword(c *gin.Context) {
	file, fileErr := c.FormFile("file")
	if fileErr != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: errors.New(invalidFileErr).Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(http.StatusBadRequest, errorResponse.NewErrorResponse())
		return
	}

	currentUserID := helpers.GetCurrentUserID(c)
	keywordService := keyword_service.KeywordService{CurrentUserID: currentUserID}

	validateTypeErr := keywordService.ValidateFileType(file.Header["Content-Type"][0])
	if validateTypeErr != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: validateTypeErr.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	filename := keywordService.UploadFile(c, file)

	parsedKeywordList, readFileErr := keywordService.ReadFile(filename)
	if readFileErr != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: readFileErr.Error(),
			Status: http.StatusUnprocessableEntity,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	// Validate if CSV has row between 1 and 1,000
	validateLengthErr := keywordService.ValidateCSVLength(len(parsedKeywordList))
	if validateLengthErr != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: validateLengthErr.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	// Save keywords to database
	saveKeywordsErr := keywordService.Save(parsedKeywordList)
	if saveKeywordsErr != nil {
		errorResponse := &api_helper.ErrorResponseObject{
			Detail: saveKeywordsErr.Error(),
			Status: http.StatusBadRequest,
		}
		c.JSON(errorResponse.Status, errorResponse.NewErrorResponse())
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
