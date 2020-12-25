package presenters

import (
	"encoding/json"
	"fmt"

	"github.com/gutakk/go-google-scraper/models"
)

type KeywordPresenter struct {
	Keyword models.Keyword
}

func (kp *KeywordPresenter) FormattedCreatedAt() string {
	return kp.Keyword.CreatedAt.Format("January 2, 2006")
}

type KeywordResult struct {
	ID                      string
	Keyword                 string
	Status                  models.KeywordStatus
	LinksCount              int
	NonAdwordsCount         int
	NonAdwordLinks          []string
	TopPositionAdwordsCount int
	TopPositionAdwordLinks  []string
	TotalAdwordsCount       int
	HtmlCode                string
	FailedReason            string
}

func (kp *KeywordPresenter) KeywordResult() KeywordResult {
	var nonAdwordLinks []string
	_ = json.Unmarshal(kp.Keyword.NonAdwordLinks, &nonAdwordLinks)

	var topPositionAdwordLinks []string
	_ = json.Unmarshal(kp.Keyword.TopPositionAdwordLinks, &topPositionAdwordLinks)

	return KeywordResult{
		ID:                      fmt.Sprint(kp.Keyword.ID),
		Keyword:                 kp.Keyword.Keyword,
		Status:                  kp.Keyword.Status,
		LinksCount:              kp.Keyword.LinksCount,
		NonAdwordsCount:         kp.Keyword.NonAdwordsCount,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: kp.Keyword.TopPositionAdwordsCount,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       kp.Keyword.TotalAdwordsCount,
		HtmlCode:                kp.Keyword.HtmlCode,
		FailedReason:            kp.Keyword.FailedReason,
	}
}
