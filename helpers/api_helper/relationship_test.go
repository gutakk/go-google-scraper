package api_helper_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/helpers/api_helper"

	"gopkg.in/go-playground/assert.v1"
)

func TestGetRelationshipsWithValidDataResponseObject(t *testing.T) {
	relationshipsResponseObject := api_helper.RelationshipsObject{
		ID:   "1",
		Type: "test",
	}

	result := relationshipsResponseObject.JSONAPIFormat()

	expected := make(map[string]api_helper.RelationshipsResponse)
	expected["test"] = api_helper.RelationshipsResponse{
		Data: api_helper.RelationshipsObject{
			ID:   "1",
			Type: "test",
		},
	}

	assert.Equal(t, expected, result)
}
