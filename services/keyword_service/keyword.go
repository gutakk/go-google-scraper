package keyword_service

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	errorHandler "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"
)

const (
	fileFormatError         = "file must be CSV format"
	fileLengthError         = "CSV file must contain between 1 to 1000 keywords"
	invalidDataError        = "invalid data"
	somethingWentWrongError = "something went wrong, please try again"
)

type KeywordService struct {
	CurrentUserID uint
}

func (k *KeywordService) GetAll() ([]models.Keyword, error) {
	condition := make(map[string]interface{})
	condition["user_id"] = k.CurrentUserID

	keywords, err := models.GetKeywordsBy(condition)
	if err != nil {
		return nil, errorHandler.DatabaseErrorMessage(err)
	}

	return keywords, nil
}

func (k *KeywordService) Save(parsedKeywordList []string) ([]models.Keyword, error) {
	// Check if record is empty slices
	if len(parsedKeywordList) == 0 {
		return nil, errors.New(invalidDataError)
	}

	var bulkData = []models.Keyword{}
	// Create bulk data
	for _, value := range parsedKeywordList {
		bulkData = append(bulkData, models.Keyword{Keyword: value, UserID: k.CurrentUserID})
	}

	savedKeywords, err := models.SaveKeywords(bulkData)
	if err != nil {
		return nil, errorHandler.DatabaseErrorMessage(err)
	}

	return savedKeywords, nil
}

func (k *KeywordService) ReadFile(filename string) ([]string, error) {
	csvfile, openErr := os.Open(filename)
	if openErr != nil {
		return nil, errors.New(somethingWentWrongError)
	}

	r := csv.NewReader(csvfile)
	var keywordList []string
	for {
		row, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.New(somethingWentWrongError)
		}
		keywordList = append(keywordList, row[0])
	}

	return keywordList, nil
}

func (k *KeywordService) UploadFile(c *gin.Context, file *multipart.FileHeader) string {
	path := "dist/"
	_ = os.Mkdir(path, 0755)
	filename := filepath.Join(path, filepath.Base(file.Filename))
	_ = c.SaveUploadedFile(file, filename)
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
