package models_test

import (
	"testing"

	"github.com/gutakk/go-google-scraper/models"

	"gopkg.in/go-playground/assert.v1"
)

func TestGetJoinedConditionsWithValidConditionsMap(t *testing.T) {
	conditions := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "testKeyword",
		},
		{
			ConditionName: "user_id",
			Value:         "testUserID",
		},
	}

	result, err := models.GetJoinedConditions(conditions)

	expected := "LOWER(keyword) LIKE LOWER('%testKeyword%') AND user_id = 'testUserID'"

	assert.Equal(t, expected, result)
	assert.Equal(t, nil, err)
}

func TestGetJoinedConditionsWithInvalidFilter(t *testing.T) {
	conditions := []models.Condition{
		{
			ConditionName: "invalid",
			Value:         "testKeyword",
		},
		{
			ConditionName: "invalid",
			Value:         "testUserID",
		},
	}

	result, err := models.GetJoinedConditions(conditions)

	assert.Equal(t, "", result)
	assert.Equal(t, "could not join conditions", err.Error())
}

func TestGetJoinedConditionsWithBlankValue(t *testing.T) {
	conditions := []models.Condition{
		{
			ConditionName: "keyword",
			Value:         "",
		},
		{
			ConditionName: "user_id",
			Value:         "",
		},
	}

	result, err := models.GetJoinedConditions(conditions)

	assert.Equal(t, "", result)
	assert.Equal(t, "could not join conditions", err.Error())
}
