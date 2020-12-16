package presenters

import (
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestFormattedCreatedAtWithValidTime(t *testing.T) {
	time := time.Date(2020, 12, 16, 12, 0, 0, 0, time.UTC)
	result := FormattedCreatedAt(time)

	assert.Equal(t, "December 16, 2020", result)
}
