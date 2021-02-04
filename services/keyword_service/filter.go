package keyword_service

import (
	"github.com/gutakk/go-google-scraper/models"
)

var FilterList = []map[string]string{
	{
		"queryString":    "filter[keyword]",
		"modelCondition": models.KeywordCondition,
	},
	{
		"queryString":    "filter[url]",
		"modelCondition": models.URLCondition,
	},
	{
		"queryString":    "filter[is_adword_advertiser]",
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
