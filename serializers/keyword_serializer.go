package serializers

import (
	"fmt"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"

	"gorm.io/datatypes"
)

type KeywordSerializer struct {
	Keyword models.Keyword
}

type KeywordJSONResponse struct {
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

func (k *KeywordSerializer) JSONAPIFormat() api_helper.DataResponse {
	if k.Keyword.Model == nil {
		return api_helper.DataResponse{}
	}

	formattedKeyword := KeywordJSONResponse{
		Keyword:                 k.Keyword.Keyword,
		Status:                  k.Keyword.Status,
		LinksCount:              k.Keyword.LinksCount,
		NonAdwordsCount:         k.Keyword.NonAdwordsCount,
		NonAdwordLinks:          k.Keyword.NonAdwordLinks,
		TopPositionAdwordsCount: k.Keyword.TopPositionAdwordsCount,
		TopPositionAdwordLinks:  k.Keyword.TopPositionAdwordLinks,
		TotalAdwordsCount:       k.Keyword.TotalAdwordsCount,
		HtmlCode:                k.Keyword.HtmlCode,
		FailedReason:            k.Keyword.FailedReason,
	}

	relationships := api_helper.RelationshipsObject{
		ID:   fmt.Sprint(k.Keyword.UserID),
		Type: models.UserType,
	}

	dataResponse := api_helper.DataResponse{
		Data: api_helper.DataResponseObject{
			ID:            fmt.Sprint(k.Keyword.ID),
			Type:          models.KeywordType,
			Attributes:    formattedKeyword,
			Relationships: relationships.JSONAPIFormat(),
		},
	}

	return dataResponse
}
