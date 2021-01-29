package filter_helper

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

func FilterValidConditions(queryString map[string][]string) []models.Condition {
	var validConditions []models.Condition

	for _, f := range FilterList {
		queryStringValue := queryString[f["queryString"]]
		if queryStringValue != nil && queryStringValue[0] != "" {
			validConditions = append(validConditions, models.Condition{
				ConditionName: f["modelCondition"],
				Value:         queryStringValue[0],
			})
		}
	}

	return validConditions
}
