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

func TestGetRelationshipsWithValidDataResponseObject(t *testing.T) {
	dataResponseObject := api_helper.DataResponseObject{
		ID:   "1",
		Type: "test",
	}

	result := dataResponseObject.GetRelationships()

	expected := make(map[string]api_helper.DataResponse)
	expected["test"] = api_helper.DataResponse{
		Data: api_helper.DataResponseObject{
			ID:            "1",
			Type:          "test",
			Attributes:    nil,
			Relationships: nil,
		},
	}

	assert.Equal(t, expected, result)
}
