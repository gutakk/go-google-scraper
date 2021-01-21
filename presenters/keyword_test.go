package presenters

import (
	"encoding/json"
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

func TestKeywordResultWithValidKeywordModels(t *testing.T) {
	nonAdwordLinks, _ := json.Marshal([]string{"non-ad-link1", "non-ad-link2"})
	topPositionAdwordLinks, _ := json.Marshal([]string{"top-ad-link1", "top-ad-link2"})

	keyword := models.Keyword{
		Keyword:                 "test-keyword",
		Status:                  models.Pending,
		LinksCount:              10,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: 10,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       10,
		HtmlCode:                "test-html",
	}

	keywordPresenter := KeywordPresenter{Keyword: keyword}
	result := keywordPresenter.KeywordLinks()

	assert.Equal(t, "test-keyword", keywordPresenter.Keyword.Keyword)
	assert.Equal(t, "pending", string(keywordPresenter.Keyword.Status))
	assert.Equal(t, 10, keywordPresenter.Keyword.LinksCount)
	assert.Equal(t, []string{"non-ad-link1", "non-ad-link2"}, result.NonAdwordLinks)
	assert.Equal(t, 10, keywordPresenter.Keyword.TopPositionAdwordsCount)
	assert.Equal(t, []string{"top-ad-link1", "top-ad-link2"}, result.TopPositionAdwordLinks)
	assert.Equal(t, 10, keywordPresenter.Keyword.TotalAdwordsCount)
	assert.Equal(t, "test-html", keywordPresenter.Keyword.HtmlCode)
}
