package keyword_service

import (
	"fmt"
)

const (
	keywordTitleQuery = "filter[keyword]"
	keywordTitleDBCol = "keyword"
)

func GetKeywordConditionsFromQueryStrings(queryStrings map[string][]string) []string {
	var conditions []string

	title := ensureQueryKeyExist(queryStrings, keywordTitleQuery)
	if title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE LOWER('%%%s%%')", keywordTitleDBCol, title))
	}

	return conditions
}

func ensureQueryKeyExist(queryStrings map[string][]string, key string) string {
	v, found := queryStrings[key]
	if found {
		return v[0]
	}
	return ""
}
