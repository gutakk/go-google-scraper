package models

import (
	"errors"
	"fmt"
	"strings"
)

type Condition struct {
	ConditionName string
	Value         string
}

// Map query string filter to where statement
var ConditionMapper = map[string]string{
	"keyword": "LOWER(keyword) LIKE LOWER('%%%s%%')",
	"user_id": "user_id = '%s'",
}

func GetJoinedConditions(conditions []Condition) (string, error) {
	var formattedConditions []string

	for _, c := range conditions {
		conditionName := c.ConditionName
		conditionValue := c.Value
		whereStatement := ConditionMapper[conditionName]

		if conditionValue != "" && whereStatement != "" {
			formattedConditions = append(formattedConditions, fmt.Sprintf(whereStatement, conditionValue))
		} else {
			return "", errors.New(couldNotJoinConditionError)
		}
	}

	joinedConditions := strings.Join(formattedConditions, " AND ")

	return joinedConditions, nil
}
