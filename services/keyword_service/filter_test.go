package keyword_service_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/services/keyword_service"

	"gopkg.in/go-playground/assert.v1"
)

func TestGetKeywordConditionsFromQueryStringsWithValidKeywordTitle(t *testing.T) {
	queryString := map[string][]string{
		"filter[keyword]": {"test"},
	}
	result := keyword_service.GetKeywordConditionsFromQueryStrings(queryString)
	expectedResult := []string{"LOWER(keyword) LIKE LOWER('%test%')"}

	assert.Equal(t, expectedResult, result)
}

func TestGetKeywordConditionsFromQueryStringsWithInvalidKey(t *testing.T) {
	queryString := map[string][]string{
		"invalid-key": {"test"},
	}
	result := keyword_service.GetKeywordConditionsFromQueryStrings(queryString)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}

func TestGetKeywordConditionsFromQueryStringsWithEmptyQueryString(t *testing.T) {
	queryString := map[string][]string{}
	result := keyword_service.GetKeywordConditionsFromQueryStrings(queryString)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}

func TestGetKeywordConditionsFromQueryStringsWithNilQueryString(t *testing.T) {
	result := keyword_service.GetKeywordConditionsFromQueryStrings(nil)
	var expectedResult []string

	assert.Equal(t, expectedResult, result)
}
