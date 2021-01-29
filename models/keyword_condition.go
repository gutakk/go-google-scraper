package models

import (
	"errors"
	"fmt"
	"strings"
)

const (
	couldNotJoinConditionError = "could not join conditions"

	// Condition name
	KeywordCondition            = "keyword"
	UserIDCondition             = "user_id"
	URLCondition                = "url"
	IsAdwordAdvertiserCondition = "is_adword_advertiser"
)

type Condition struct {
	ConditionName string
	Value         string
}

// Map query string filter to where statement
var ConditionMapper = map[string]string{
	KeywordCondition:            "LOWER(keyword) LIKE LOWER('%%%s%%')",
	UserIDCondition:             "user_id = '%s'",
	URLCondition:                "(LOWER(non_adword_links::text) LIKE '%%%[1]s%%' OR LOWER(top_position_adword_links::text) LIKE '%%%[1]s%%')",
	IsAdwordAdvertiserCondition: "total_adwords_count > 0",
}

func GetJoinedConditions(conditions []Condition) (string, error) {
	var formattedConditions []string

	for _, c := range conditions {
		conditionName := c.ConditionName
		conditionValue := c.Value
		whereStatement := ConditionMapper[conditionName]

		if conditionValue != "" && whereStatement != "" {
			if conditionName == IsAdwordAdvertiserCondition {
				formattedConditions = append(formattedConditions, whereStatement)
			} else {
				formattedConditions = append(formattedConditions, fmt.Sprintf(whereStatement, conditionValue))
			}
		} else {
			return "", errors.New(couldNotJoinConditionError)
		}
	}

	joinedConditions := strings.Join(formattedConditions, " AND ")

	return joinedConditions, nil
}
