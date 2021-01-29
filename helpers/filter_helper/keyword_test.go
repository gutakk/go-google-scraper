package filter_helper_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/filter_helper"
	"github.com/gutakk/go-google-scraper/models"

	"gopkg.in/go-playground/assert.v1"
)

func TestFilterValidConditionsWithValidQueryString(t *testing.T) {
	queryString := map[string][]string{
		"filter[keyword]": {"test"},
	}

	result := filter_helper.FilterValidConditions(queryString)

	expected := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "test",
		},
	}

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithoutQueryString(t *testing.T) {
	result := filter_helper.FilterValidConditions(nil)

	var expected []models.Condition

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithInvalidQueryString(t *testing.T) {
	queryString := map[string][]string{
		"filter[invalid]": {"test"},
	}

	result := filter_helper.FilterValidConditions(queryString)

	var expected []models.Condition

	assert.Equal(t, expected, result)
}

func TestFilterValidConditionsWithBlankQueryStringValue(t *testing.T) {
	queryString := map[string][]string{
		"filter[keyword]": {""},
	}

	result := filter_helper.FilterValidConditions(queryString)

	var expected []models.Condition

	assert.Equal(t, expected, result)
}
