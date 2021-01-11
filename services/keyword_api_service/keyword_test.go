package keyword_api_service_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/keyword_api_service"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/gorm"
)

func TestJSONAPIFormatKeywordsResponseWithValidKeywords(t *testing.T) {
	keywords := keyword_api_service.KeywordsResponse{
		Keywords: []models.Keyword{
			{
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
			},
			{
				Model:                   &gorm.Model{ID: 2},
				Keyword:                 "testKeyword2",
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
				User:                    &models.User{}},
		},
	}

	result := keywords.JSONAPIFormatKeywordsResponse()

	relationships := make(map[string]api_helper.DataResponse)
	relationships["user"] = api_helper.DataResponse{
		Data: api_helper.DataResponseObject{
			ID:   "1",
			Type: "user",
		},
	}

	expected := api_helper.DataResponseArray{
		Data: []api_helper.DataResponseObject{
			{
				ID:   "1",
				Type: "keyword",
				Attributes: models.Keyword{
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
			{
				ID:   "2",
				Type: "keyword",
				Attributes: models.Keyword{
					Keyword:                 "testKeyword2",
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
		},
	}

	assert.Equal(t, expected, result)
}

func TestJSONAPIFormatKeywordsResponseWithBlankKeywords(t *testing.T) {
	keywords := keyword_api_service.KeywordsResponse{}
	result := keywords.JSONAPIFormatKeywordsResponse()

	assert.Equal(t, api_helper.DataResponseArray{}, result)
}
