package presenters

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/gutakk/go-google-scraper/models"

	"github.com/bxcodec/faker/v3"
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

func TestKeywordResultWithValidKeywordModels(t *testing.T) {
	nonAdwordLinks, _ := json.Marshal([]string{"non-ad-link1", "non-ad-link2"})
	topPositionAdwordLinks, _ := json.Marshal([]string{"top-ad-link1", "top-ad-link2"})

	keyword := models.Keyword{
		Keyword:                 faker.Name(),
		Status:                  models.Pending,
		LinksCount:              10,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: 10,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       10,
		HtmlCode:                "test-html",
	}

	keywordPresenter := KeywordPresenter{Keyword: keyword}
	result := keywordPresenter.KeywordResult()

	assert.Equal(t, "0", result.ID)
	assert.Equal(t, keyword.Keyword, result.Keyword)
	assert.Equal(t, models.Pending, result.Status)
	assert.Equal(t, 10, result.LinksCount)
	assert.Equal(t, 10, result.TopPositionAdwordsCount)
	assert.Equal(t, []string{"non-ad-link1", "non-ad-link2"}, result.NonAdwordLinks)
	assert.Equal(t, 10, result.TotalAdwordsCount)
	assert.Equal(t, []string{"top-ad-link1", "top-ad-link2"}, result.TopPositionAdwordLinks)
	assert.Equal(t, "test-html", result.HtmlCode)
	assert.Equal(t, "", result.FailedReason)
}
