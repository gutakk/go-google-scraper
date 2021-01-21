package view

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestSearchResultPartialWithValidParams(t *testing.T) {
	result := SearchResultPartial("test-title", "test-value")

	assert.Equal(t, "test-title", result["title"])
	assert.Equal(t, "test-value", result["value"])
}
