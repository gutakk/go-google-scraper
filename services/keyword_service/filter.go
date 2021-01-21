package keyword_service

import (
	"fmt"
)

const (
	keywordTitleQuery = "keyword-title"
	keywordTitleDBCol = "keyword"
)

func GetConditionFromQuery(queryStrings map[string][]string) []string {
	var conditions []string

	title := keywordTitle(queryStrings)
	if title != "" {
		conditions = append(conditions, fmt.Sprintf("LOWER(%s) LIKE LOWER('%%%s%%')", keywordTitleDBCol, title))
	}

	return conditions
}

func keywordTitle(queryStrings map[string][]string) string {
	v, found := queryStrings[keywordTitleQuery]
	if found {
		return v[0]
	}
	return ""
}
