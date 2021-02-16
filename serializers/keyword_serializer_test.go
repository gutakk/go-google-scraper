package serializers_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/serializers"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

func testJsonAPIFormatKeywordResponseWithValidKeyword(t *testing.T) {
	keyword := models.Keyword{
		Model:                   &gorm.Model{ID: 1},
		Keyword:                 "testKeyword1",
		Status:                  models.Pending,
		LinksCount:              1,
		NonAdwordsCount:         1,
		NonAdwordLinks:          []byte("testNonAdwordLinks"),
		TopPositionAdwordsCount: 1,
		TopPositionAdwordLinks:  []byte("testTopPositionAdwordLinks"),
		TotalAdwordsCount:       1,
		UserID:                  1,
		HtmlCode:                "testHTML",
		FailedReason:            "",
		User:                    &models.User{},
	}

	keywordSerializer := serializers.KeywordSerializer{Keyword: keyword}
	result := keywordSerializer.JSONAPIFormat()

	relationships := make(map[string]api_helper.RelationshipsResponse)
	relationships["user"] = api_helper.RelationshipsResponse{
		Data: api_helper.RelationshipsObject{
			ID:   "1",
			Type: "user",
		},
	}

	expected := api_helper.DataResponse{
		Data: api_helper.DataResponseObject{
			ID:   "1",
			Type: "keyword",
			Attributes: serializers.KeywordJSONResponse{
				Keyword:                 "testKeyword1",
				Status:                  models.Pending,
				LinksCount:              1,
				NonAdwordsCount:         1,
				NonAdwordLinks:          []byte("testNonAdwordLinks"),
				TopPositionAdwordsCount: 1,
				TopPositionAdwordLinks:  []byte("testTopPositionAdwordLinks"),
				TotalAdwordsCount:       1,
				HtmlCode:                "testHTML",
				FailedReason:            "",
			},
			Relationships: relationships,
		},
	}

	assert.Equal(t, expected, result)
}

func testJsonAPIFormatKeywordResponseWithBlankKeyword(t *testing.T) {
	keyword := models.Keyword{}
	keywordSerializer := serializers.KeywordSerializer{Keyword: keyword}
	result := keywordSerializer.JSONAPIFormat()

	assert.Equal(t, api_helper.DataResponse{Data: api_helper.DataResponseObject{}}, result)
}
