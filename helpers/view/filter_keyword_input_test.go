package view_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/view"

	"gopkg.in/go-playground/assert.v1"
)

func TestFilterKeywordPartialInputWithValidTextQueryString(t *testing.T) {
	filter := map[string][]string{
		"filter[text]": {"Lorem Ipsum"},
	}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "text", "testClassName", "filter[text]", "testPlaceholder")

	assert.Equal(t, "Lorem Ipsum", result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "text", result["inputType"])
	assert.Equal(t, "testClassName", result["className"])
	assert.Equal(t, "filter[text]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}

func TestFilterKeywordPartialInputWithValidCheckboxQueryString(t *testing.T) {
	filter := map[string][]string{
		"filter[checkbox]": {"true"},
	}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "checkbox", "testClassName", "filter[checkbox]", "testPlaceholder")

	assert.Equal(t, true, result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "checkbox", result["inputType"])
	assert.Equal(t, "testClassName", result["className"])
	assert.Equal(t, "filter[checkbox]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}

func TestFilterKeywordPartialInputWithoutCheckboxQueryString(t *testing.T) {
	filter := map[string][]string{
		"filter[invalid]": {"test"},
	}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "checkbox", "testClassName", "filter[checkbox]", "testPlaceholder")

	assert.Equal(t, false, result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "checkbox", result["inputType"])
	assert.Equal(t, "testClassName", result["className"])
	assert.Equal(t, "filter[checkbox]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}

func TestFilterKeywordPartialInputWithInvalidCheckboxQueryString(t *testing.T) {
	filter := map[string][]string{
		"filter[checkbox]": {"invalid"},
	}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "checkbox", "testClassName", "filter[test]", "testPlaceholder")

	assert.Equal(t, false, result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "checkbox", result["inputType"])
	assert.Equal(t, "testClassName", result["className"])
	assert.Equal(t, "filter[test]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}

func TestFilterKeywordPartialInputWithBlankFilter(t *testing.T) {
	filter := map[string][]string{}
	result := view.FilterKeywordPartialInput(filter, "testLabel", "checkbox", "testClassName", "filter[test]", "testPlaceholder")

	assert.Equal(t, nil, result["filter"])
	assert.Equal(t, "testLabel", result["filterLabel"])
	assert.Equal(t, "checkbox", result["inputType"])
	assert.Equal(t, "testClassName", result["className"])
	assert.Equal(t, "filter[test]", result["filterName"])
	assert.Equal(t, "testPlaceholder", result["filterPlaceholder"])
}
