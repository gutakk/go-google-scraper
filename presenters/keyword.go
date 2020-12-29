package presenters

import (
	"encoding/json"

	"github.com/gutakk/go-google-scraper/models"
)

type KeywordPresenter struct {
	Keyword models.Keyword
}

type KeywordResult struct {
	NonAdwordLinks         []string
	TopPositionAdwordLinks []string
}

func (kp *KeywordPresenter) FormattedCreatedAt() string {
	return kp.Keyword.CreatedAt.Format("January 2, 2006")
}

func (kp *KeywordPresenter) KeywordResult() KeywordResult {
	var nonAdwordLinks []string
	_ = json.Unmarshal(kp.Keyword.NonAdwordLinks, &nonAdwordLinks)

	var topPositionAdwordLinks []string
	_ = json.Unmarshal(kp.Keyword.TopPositionAdwordLinks, &topPositionAdwordLinks)

	return KeywordResult{
		NonAdwordLinks:         nonAdwordLinks,
		TopPositionAdwordLinks: topPositionAdwordLinks,
	}
}
