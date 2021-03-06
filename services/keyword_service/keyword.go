package keyword_service

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gutakk/go-google-scraper/db"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/google_search_service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const (
	cannotOpenFileError = "file cannot be opened"
	cannotReadFileError = "file cannot be read"
	fileFormatError     = "file must be CSV format"
	fileLengthError     = "CSV file must contain between 1 to 1000 keywords"
	invalidDataError    = "invalid data"
)

func (k *KeywordService) GetKeywords(conditions []models.Condition) ([]models.Keyword, error) {
	conditions = append(conditions, models.Condition{
		ConditionName: models.UserIDCondition,
		Value:         fmt.Sprint(k.CurrentUserID),
	})

	keywords, err := models.GetKeywordsBy(conditions)
	if err != nil {
		return nil, errorHandler.DatabaseErrorMessage(err)
	}

	return keywords, nil
}

func (k *KeywordService) GetKeywordResult(keywordID interface{}) (models.Keyword, error) {
	condition := make(map[string]interface{})
	condition["id"] = keywordID
	condition["user_id"] = k.CurrentUserID

	keyword, err := models.GetKeywordBy(condition)
	if err != nil {
		return models.Keyword{}, errorHandler.DatabaseErrorMessage(err)
	}

	return keyword, nil
}

func (k *KeywordService) Save(parsedKeywordList []string) error {
	// Check if record is empty slices
	if len(parsedKeywordList) == 0 {
		return errors.New(invalidDataError)
	}

	for _, value := range parsedKeywordList {
		keyword := models.Keyword{Keyword: value, UserID: k.CurrentUserID}

		err := db.GetDB().Transaction(func(tx *gorm.DB) error {
			savedKeyword, err := models.SaveKeyword(keyword, tx)
			if err != nil {
				return errorHandler.DatabaseErrorMessage(err)
			}

			return google_search_service.EnqueueSearchJob(savedKeyword)
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (k *KeywordService) ReadFile(filename string) ([]string, error) {
	csvfile, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(cannotOpenFileError)
	}

	r := csv.NewReader(csvfile)
	var keywordList []string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New(cannotReadFileError)
		}
		keywordList = append(keywordList, row[0])
	}

	return keywordList, nil
}

func (k *KeywordService) UploadFile(c *gin.Context, file *multipart.FileHeader) string {
	path := "dist/"
	err := os.Mkdir(path, 0755)
	if err != nil {
		log.Error("Failed to create directory: ", err)
	}

	filename := filepath.Join(path, filepath.Base(file.Filename))
	err = c.SaveUploadedFile(file, filename)
	if err != nil {
		log.Error("Failed to save uploaded file: ", err)
	}

	return filename
}

func (k *KeywordService) ValidateCSVLength(row int) error {
	if row <= 0 || row > 1000 {
		return errors.New(fileLengthError)
	}
	return nil
}

func (k *KeywordService) ValidateFileType(fileType string) error {
	if fileType != "text/csv" {
		return errors.New(fileFormatError)
	}
	return nil
}
