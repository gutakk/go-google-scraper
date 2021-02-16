package serializers

import (
	"fmt"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"

	"gorm.io/datatypes"
)

type KeywordsSerializer struct {
	Keywords []models.Keyword
}

type KeywordsJSONResponse struct {
	Keyword                 string               `json:"keyword"`
	Status                  models.KeywordStatus `json:"status"`
	LinksCount              int                  `json:"links_count"`
	NonAdwordsCount         int                  `json:"non_adwords_count"`
	NonAdwordLinks          datatypes.JSON       `json:"non_adword_links"`
	TopPositionAdwordsCount int                  `json:"top_position_adwords_count"`
	TopPositionAdwordLinks  datatypes.JSON       `json:"top_position_adword_links"`
	TotalAdwordsCount       int                  `json:"total_adwords_count"`
	HtmlCode                string               `json:"html_code"`
	FailedReason            string               `json:"failed_reason"`
}

func (k *KeywordsSerializer) JSONAPIFormat() api_helper.DataResponseList {
	formattedKeywords := api_helper.DataResponseList{}
	formattedKeywords.Data = []api_helper.DataResponseObject{}

	for _, value := range k.Keywords {
		keywordResp := KeywordsJSONResponse{
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

		relationships := api_helper.RelationshipsObject{
			ID:   fmt.Sprint(value.UserID),
			Type: models.UserType,
		}

		dataResponseObject := api_helper.DataResponseObject{
			ID:            fmt.Sprint(value.ID),
			Type:          models.KeywordType,
			Attributes:    keywordResp,
			Relationships: relationships.JSONAPIFormat(),
		}

		formattedKeywords.Data = append(formattedKeywords.Data, dataResponseObject)
	}

	return formattedKeywords
}
