package presenters

import (
	"encoding/json"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
)

type KeywordPresenter struct {
	Keyword models.Keyword
}

type KeywordLinks struct {
	NonAdwordLinks         []string
	TopPositionAdwordLinks []string
}

func (kp *KeywordPresenter) FormattedCreatedAt() string {
	return kp.Keyword.CreatedAt.Format("January 2, 2006")
}

func (kp *KeywordPresenter) KeywordLinks() KeywordLinks {
	var nonAdwordLinks []string
	err := json.Unmarshal(kp.Keyword.NonAdwordLinks, &nonAdwordLinks)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	var topPositionAdwordLinks []string
	err = json.Unmarshal(kp.Keyword.TopPositionAdwordLinks, &topPositionAdwordLinks)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	return KeywordLinks{
		NonAdwordLinks:         nonAdwordLinks,
		TopPositionAdwordLinks: topPositionAdwordLinks,
	}
}
