package models

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestValidateFileTypeWithValidFileType(t *testing.T) {
	result := ValidateFileType("text/csv")

	assert.Equal(t, nil, result)
}

func TestValidateFileTypeWithInvalidFileType(t *testing.T) {
	result := ValidateFileType("test")

	assert.Equal(t, "File must be CSV format", result.Error())
}

func TestValidateCSVLengthWithOneRow(t *testing.T) {
	result := ValidateCSVLength(1)

	assert.Equal(t, nil, result)
}

func TestValidateCSVLengthWithOneThoudsandRow(t *testing.T) {
	result := ValidateCSVLength(1000)

	assert.Equal(t, nil, result)
}

func TestValidateCSVLengthWithZeroRow(t *testing.T) {
	result := ValidateCSVLength(0)

	assert.Equal(t, "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func TestValidateCSVLengthWithMinusOneRow(t *testing.T) {
	result := ValidateCSVLength(-1)

	assert.Equal(t, "CSV file must contain between 1 to 1000 keywords", result.Error())
}

func TestValidateCSVLengthWithOneThoudsandOneRow(t *testing.T) {
	result := ValidateCSVLength(1001)

	assert.Equal(t, "CSV file must contain between 1 to 1000 keywords", result.Error())
}
