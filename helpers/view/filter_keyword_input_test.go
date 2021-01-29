package view_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/view"

	"gopkg.in/go-playground/assert.v1"
)

func TestFilterKeywordPartialInputWithValidParams(t *testing.T) {
	filter := map[string][]string{
		"filter[test]": {"Lorem Ipsum"},
	}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "filter[test]", "testPlaceholder")

	assert.Equal(t, "Lorem Ipsum", result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "filter[test]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}
