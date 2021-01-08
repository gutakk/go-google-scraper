package string_helper_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/string_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestCamelCaseToSnakeCaseWithValidCamelString(t *testing.T) {
	result := string_helper.CamelCaseToSnakeCase("helloWorld")

	assert.Equal(t, "hello_world", result)
}

func TestCamelCaseToSnakeCaseWithValidCamelCapitalizeString(t *testing.T) {
	result := string_helper.CamelCaseToSnakeCase("HelloWorld")

	assert.Equal(t, "hello_world", result)
}

func TestCamelCaseToSnakeCaseWithAllUpperString(t *testing.T) {
	result := string_helper.CamelCaseToSnakeCase("HELLOWORLD")

	assert.Equal(t, "helloworld", result)
}

func TestCamelCaseToSnakeCaseWithMixingStringAndNumeric(t *testing.T) {
	result := string_helper.CamelCaseToSnakeCase("Hello16World")

	assert.Equal(t, "hello_16_world", result)
}

func TestCamelCaseToSnakeCaseWithMorethanOneUpperCharacter(t *testing.T) {
	result := string_helper.CamelCaseToSnakeCase("helloJSONWorld")

	assert.Equal(t, "hello_json_world", result)
}
