package keyword_service

import (
	"fmt"

	"github.com/gutakk/go-google-scraper/models"
)

var FilterList = []map[string]string{
	{
		"queryString":    fmt.Sprintf("filter[%s]", models.KeywordCondition),
		"modelCondition": models.KeywordCondition,
	},
	{
		"queryString":    fmt.Sprintf("filter[%s]", models.URLCondition),
		"modelCondition": models.URLCondition,
	},
	{
		"queryString":    fmt.Sprintf("filter[%s]", models.IsAdwordAdvertiserCondition),
		"modelCondition": models.IsAdwordAdvertiserCondition,
	},
}

func (k *KeywordService) FilterValidConditions() []models.Condition {
	var validConditions []models.Condition

	for _, f := range FilterList {
		queryStringValue := k.QueryString[f["queryString"]]
		if queryStringValue != nil && queryStringValue[0] != "" {
			validConditions = append(validConditions, models.Condition{
				ConditionName: f["modelCondition"],
				Value:         queryStringValue[0],
			})
		}
	}

	return validConditions
}
