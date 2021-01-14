package view_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/view"

	"gopkg.in/go-playground/assert.v1"
)

func TestFilterKeywordPartialInputWithValidParams(t *testing.T) {
	query := map[string][]string{
		"testQuery": {"Lorem Ipsum"},
	}
	result := view.FilterKeywordPartialInput(query, "testLabel", "testName", "testPlaceholder")

	assert.Equal(t, map[string][]string{"testQuery": {"Lorem Ipsum"}}, result["query"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "testName", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}
