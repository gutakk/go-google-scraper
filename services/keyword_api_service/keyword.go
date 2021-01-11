package keyword_api_service

import (
	"fmt"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"
)

type KeywordsResponse struct {
	Keywords []models.Keyword
}

func (k *KeywordsResponse) JSONAPIFormatKeywordsResponse() api_helper.DataResponseArray {
	formattedKeywords := api_helper.DataResponseArray{}

	for _, value := range k.Keywords {
		keyword := models.Keyword{
			Keyword:                 value.Keyword,
			Status:                  value.Status,
			LinksCount:              value.LinksCount,
			NonAdwordsCount:         value.NonAdwordsCount,
			NonAdwordLinks:          value.NonAdwordLinks,
			TopPositionAdwordsCount: value.TopPositionAdwordsCount,
			TopPositionAdwordLinks:  value.TopPositionAdwordLinks,
			TotalAdwordsCount:       value.TotalAdwordsCount,
			HtmlCode:                value.HtmlCode,
			FailedReason:            value.FailedReason,
		}

		relationships := api_helper.DataResponseObject{
			ID:   fmt.Sprint(value.UserID),
			Type: "user",
		}

		dataResponseObject := api_helper.DataResponseObject{
			ID:            fmt.Sprint(value.ID),
			Type:          "keyword",
			Attributes:    keyword,
			Relationships: relationships.GetRelationships(),
		}

		formattedKeywords.Data = append(formattedKeywords.Data, dataResponseObject)
	}

	return formattedKeywords
}
