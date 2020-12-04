package keyword_service

import (
	"errors"

	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	invalidDataError = "Invalid data"
)

type Keyword struct {
	CurrentUserID uint
}

func (k *Keyword) Save(record []string) ([]models.Keyword, error) {
	// Check if record is empty slices
	if len(record) == 0 {
		return nil, errors.New(invalidDataError)
	}

	var bulkData = []models.Keyword{}
	// Create bulk data
	for _, value := range record {
		bulkData = append(bulkData, models.Keyword{Keyword: value, UserID: k.CurrentUserID})
	}

	keywords, err := models.SaveKeywords(bulkData)
	if err != nil {
		return nil, errorHandler.DatabaseErrorMessage(err)
	}

	return keywords, nil
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
