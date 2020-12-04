package keyword_service

import (
	"github.com/gutakk/go-google-scraper/models"
)

type Keyword struct {
	CurrentUserID uint
}

func (k *Keyword) GetAll() ([]models.Keyword, error) {
	condition := make(map[string]interface{})
	condition["user_id"] = k.CurrentUserID

	keywords, err := models.GetKeywords(condition)
	if err != nil {
		return nil, err
	}

	return keywords, nil
}
