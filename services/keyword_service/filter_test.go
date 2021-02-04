package keyword_service_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/keyword_service"

	"gopkg.in/go-playground/assert.v1"
)

func TestFilterValidConditionsWithValidQueryStrings(t *testing.T) {
	queryString := map[string][]string{
		"filter[keyword]":              {"testKeyword"},
		"filter[url]":                  {"testURL"},
		"filter[is_adword_advertiser]": {"testAdword"},
	}

	keywordService := keyword_service.KeywordService{QueryString: queryString}
	result := keywordService.FilterValidConditions()

	expected := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "testKeyword",
		},
		{
			ConditionName: "url",
			Value:         "testURL",
		},
		{
			ConditionName: "isAdwordAdvertiser",
			Value:         "testAdword",
		},
	}

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithoutQueryString(t *testing.T) {
	keywordService := keyword_service.KeywordService{}
	result := keywordService.FilterValidConditions()

	var expected []models.Condition

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithInvalidQueryString(t *testing.T) {
	queryString := map[string][]string{
		"filter[invalid]": {"test"},
	}

	keywordService := keyword_service.KeywordService{QueryString: queryString}
	result := keywordService.FilterValidConditions()

	var expected []models.Condition

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithBlankQueryStringValue(t *testing.T) {
	queryString := map[string][]string{
		"filter[keyword]": {""},
	}

	keywordService := keyword_service.KeywordService{QueryString: queryString}
	result := keywordService.FilterValidConditions()

	var expected []models.Condition

	assert.Equal(t, expected, result)
}
