package google_scraping_service

import (
	"encoding/json"

	"github.com/gutakk/go-google-scraper/models"
)

func UpdateKeywordStatus(keywordID uint, status models.KeywordStatus) error {
	err := models.UpdateKeywordByID(keywordID, models.Keyword{Status: status})
	if err != nil {
		return err
	}

	return nil
}

func UpdateKeywordWithParsingResult(keywordID uint, parsingResult ParsingResult) error {
	nonAdwordLinks, nonAdwordLinksParsingErr := json.Marshal(parsingResult.NonAdwordLinks)
	if nonAdwordLinksParsingErr != nil {
		return nonAdwordLinksParsingErr
	}

	topPositionAdwordLinks, topPositionAdwordLinksParsingErr := json.Marshal(parsingResult.TopPositionAdwordLinks)
	if topPositionAdwordLinksParsingErr != nil {
		return topPositionAdwordLinksParsingErr
	}

	newKeyword := models.Keyword{
		Status:                  models.Processed,
		LinksCount:              parsingResult.LinksCount,
		NonAdwordsCount:         parsingResult.NonAdwordsCount,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: parsingResult.TopPostionAdwordsCount,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       parsingResult.TotalAdwordsCount,
		HtmlCode:                parsingResult.HtmlCode,
	}

	err := models.UpdateKeywordByID(keywordID, newKeyword)
	if err != nil {
		return err
	}

	return nil
}
