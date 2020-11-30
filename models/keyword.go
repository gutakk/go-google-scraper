package models

import (
	"errors"

	"gorm.io/gorm"
)

const (
	blankFileError  = "CSV file can't be blank"
	fileFormatError = "File must be CSV format"
	fileLengthError = "CSV file must contain between 1 to 1000 keywords"
)

type Keyword struct {
	gorm.Model
	Keyword string `gorm:"notNull;index"`
}

func ValidateFileFormat(fileType string) error {
	if fileType != "text/csv" {
		return errors.New(fileFormatError)
	}
	return nil
}

func ValidateCSVLength(row int) error {
	if row <= 0 {
		return errors.New(blankFileError)
	} else if row > 1000 {
		return errors.New(fileLengthError)
	}
	return nil
}
