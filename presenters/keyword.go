package presenters

import (
	"github.com/gutakk/go-google-scraper/models"
)

type KeywordPresenter struct {
	Keyword models.Keyword
}

func (kp *KeywordPresenter) FormattedCreatedAt() string {
	return kp.Keyword.CreatedAt.Format("January 2, 2006")
}
