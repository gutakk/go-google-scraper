package presenters

import (
	"testing"
	"time"

	"github.com/gutakk/go-google-scraper/models"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

func TestFormattedCreatedAtWithValidTime(t *testing.T) {
	time := time.Date(2020, 12, 16, 12, 0, 0, 0, time.UTC)
	keywordPresenter := KeywordPresenter{
		Keyword: models.Keyword{
			Model: gorm.Model{CreatedAt: time},
		},
	}

	result := keywordPresenter.FormattedCreatedAt()

	assert.Equal(t, "December 16, 2020", result)
}
