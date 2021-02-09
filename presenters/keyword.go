package presenters

import (
	"encoding/json"

	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"

	log "github.com/sirupsen/logrus"
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
		log.Error(errorHelper.JSONUnmarshalFailure, err)
	}

	var topPositionAdwordLinks []string
	err = json.Unmarshal(kp.Keyword.TopPositionAdwordLinks, &topPositionAdwordLinks)
	if err != nil {
		log.Error(errorHelper.JSONUnmarshalFailure, err)
	}

	return KeywordLinks{
		NonAdwordLinks:         nonAdwordLinks,
		TopPositionAdwordLinks: topPositionAdwordLinks,
	}
}
