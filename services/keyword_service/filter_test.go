package keyword_service_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/services/keyword_service"

	"gopkg.in/go-playground/assert.v1"
)

func TestGetConditionFromQueryWithValidQueryStringKeywordTitle(t *testing.T) {
	queryString := map[string][]string{
		"keyword-title": {"test"},
	}
	result := keyword_service.GetConditionFromQuery(queryString)
	expectedResult := []string{"LOWER(keyword) LIKE LOWER('%test%')"}

	assert.Equal(t, expectedResult, result)
}

func TestGetConditionFromQueryWithInvalidQueryStringKey(t *testing.T) {
	queryString := map[string][]string{
		"invalid-key": {"test"},
	}
	result := keyword_service.GetConditionFromQuery(queryString)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}

func TestGetConditionFromQueryWithEmptyQueryString(t *testing.T) {
	queryString := map[string][]string{}
	result := keyword_service.GetConditionFromQuery(queryString)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}

func TestGetConditionFromQueryWithNilQueryString(t *testing.T) {
	result := keyword_service.GetConditionFromQuery(nil)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}
