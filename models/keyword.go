package models

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"gorm.io/gorm"
)

const (
	fileFormatError         = "file must be CSV format"
	fileLengthError         = "CSV file must contain between 1 to 1000 keywords"
	invalidDataError        = "invalid data"
	somethingWentWrongError = "something went wrong, please try again"
)

type Keyword struct {
	gorm.Model
	Keyword string `gorm:"notNull;index"`
	UserID  uint
	User    User
}

func UploadFile(c *gin.Context, file *multipart.FileHeader) string {
	path := "dist/"
	_ = os.Mkdir(path, 0755)
	filename := filepath.Join(path, filepath.Base(file.Filename))
	_ = c.SaveUploadedFile(file, filename)
	return filename
}

func ReadFile(filename string) ([]string, error) {
	csvfile, openErr := os.Open(filename)
	if openErr != nil {
		return nil, errors.New(somethingWentWrongError)
	}

	r := csv.NewReader(csvfile)
	var record []string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New(somethingWentWrongError)
		}
		record = append(record, row[0])
	}

	return record, nil
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

func SaveKeywords(userID uint, record []string) ([]Keyword, error) {
	// Check if record is empty slices
	if len(record) == 0 {
		return nil, errors.New(invalidDataError)
	}

	var keywords = []Keyword{}
	// Create bulk data
	for _, value := range record {
		keywords = append(keywords, Keyword{Keyword: value, UserID: userID})
	}

	// Insert bulk data
	if result := db.GetDB().Create(&keywords); result.Error != nil {
		return nil, errorHandler.DatabaseErrorMessage(result.Error)
	}

	return keywords, nil
}
