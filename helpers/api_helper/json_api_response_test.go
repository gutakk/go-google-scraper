package api_helper_test

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestNewErrorResponseWithValidErrorResponseObject(t *testing.T) {
	errorResponseObject := &api_helper.ErrorResponseObject{
		Detail: "test-detail",
		Status: 999,
	}
	errorResponse := errorResponseObject.NewErrorResponse()
	expectedResult := gin.H{
		"errors": []api_helper.ErrorResponseObject{{
			Detail: "test-detail",
			Status: 999,
		}},
	}

	assert.Equal(t, expectedResult, errorResponse)
}

func TestNewErrorResponseWithMissingSomeFieldOnErrorResponseObject(t *testing.T) {
	errorResponseObject := &api_helper.ErrorResponseObject{
		Detail: "test-detail",
		Status: 999,
	}
	errorResponse := errorResponseObject.NewErrorResponse()
	expectedResult := gin.H{
		"errors": []api_helper.ErrorResponseObject{{
			Detail: "test-detail",
			Status: 999,
		}},
	}

	assert.Equal(t, expectedResult, errorResponse)
}
