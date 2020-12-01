package models

import (
	"errors"

	"github.com/gutakk/go-google-scraper/db"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"gorm.io/gorm"
)

const (
	fileFormatError  = "File must be CSV format"
	fileLengthError  = "CSV file must contain between 1 to 1000 keywords"
	invalidDataError = "Invalid data"
)

type Keyword struct {
	gorm.Model
	Keyword string `gorm:"notNull;index"`
}

func ValidateFileType(fileType string) error {
	if fileType != "text/csv" {
		return errors.New(fileFormatError)
	}
	return nil
}

func ValidateCSVLength(row int) error {
	if row <= 0 || row > 1000 {
		return errors.New(fileLengthError)
	}
	return nil
}

func SaveKeywords(record [][]string) ([]Keyword, error) {
	var keywords = []Keyword{}

	// Check if record more than 0
	if len(record) <= 0 {
		return nil, errors.New(invalidDataError)
	}

	// Create bulk data
	for _, v := range record {
		// Check if nested slice is empty slice
		if len(v) <= 0 {
			return nil, errors.New(invalidDataError)
		}
		keywords = append(keywords, Keyword{Keyword: v[0]})
	}

	// Insert bulk data
	if result := db.GetDB().Create(&keywords); result.Error != nil {
		return nil, errorHandler.DatabaseErrorMessage(result.Error)
	}

	return keywords, nil
}
