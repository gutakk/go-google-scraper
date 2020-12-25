package google_search_service

import (
	"encoding/json"

	"github.com/gutakk/go-google-scraper/models"
)

func UpdateKeywordStatus(keywordID uint, status models.KeywordStatus, failedReason error) error {
	var keyword models.Keyword
	keyword.Status = status

	if failedReason != nil {
		keyword.FailedReason = failedReason.Error()
	}

	err := models.UpdateKeyword(keywordID, keyword)
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

	keyword := models.Keyword{
		Status:                  models.Processed,
		LinksCount:              parsingResult.LinksCount,
		NonAdwordsCount:         parsingResult.NonAdwordsCount,
		NonAdwordLinks:          nonAdwordLinks,
		TopPositionAdwordsCount: parsingResult.TopPostionAdwordsCount,
		TopPositionAdwordLinks:  topPositionAdwordLinks,
		TotalAdwordsCount:       parsingResult.TotalAdwordsCount,
		HtmlCode:                parsingResult.HtmlCode,
	}

	err := models.UpdateKeyword(keywordID, keyword)
	if err != nil {
		return err
	}

	return nil
}
