package api_helper_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestConstructErrorResponseWithValidErrorResponseObject(t *testing.T) {
	errorResponseObject := api_helper.ErrorResponseObject{
		Title:  "test-error",
		Detail: "test-detail",
		Status: 999,
	}
	errorResponse := errorResponseObject.ConstructErrorResponse()
	expectedResult := gin.H{
		"errors": []api_helper.ErrorResponseObject{{
			Title:  "test-error",
			Detail: "test-detail",
			Status: 999,
		}},
	}

	assert.Equal(t, expectedResult, errorResponse)
}

func TestConstructErrorResponseWithMissingSomeFieldOnErrorResponseObject(t *testing.T) {
	errorResponseObject := api_helper.ErrorResponseObject{
		Detail: "test-detail",
		Status: 999,
	}
	errorResponse := errorResponseObject.ConstructErrorResponse()
	expectedResult := gin.H{
		"errors": []api_helper.ErrorResponseObject{{
			Detail: "test-detail",
			Status: 999,
		}},
	}

	assert.Equal(t, expectedResult, errorResponse)
}
