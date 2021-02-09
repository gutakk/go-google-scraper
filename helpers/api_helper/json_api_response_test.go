package api_helper_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"

	"github.com/gin-gonic/gin"
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
	}
	errorResponse := errorResponseObject.NewErrorResponse()
	expectedResult := gin.H{
		"errors": []api_helper.ErrorResponseObject{{
			Detail: "test-detail",
			Status: 0,
		}},
	}

	assert.Equal(t, expectedResult, errorResponse)
}
