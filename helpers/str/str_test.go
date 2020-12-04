package str

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestCapitalizeFirstWithValidString(t *testing.T) {
	result := CapitalizeFirst("hello world")

	assert.Equal(t, "Hello world", result)
}

func TestCapitalizeFirstWithEmptyString(t *testing.T) {
	result := CapitalizeFirst("")

	assert.Equal(t, "", result)
}
