package string_helper

import (
	"strings"

	"github.com/fatih/camelcase"
)

func CamelCaseToSnakeCase(s string) string {
	splitted := camelcase.Split(s)
	lowerStrings := []string{}

	for _, v := range splitted {
		lowerStrings = append(lowerStrings, strings.ToLower(v))
	}

	return strings.Join(lowerStrings, "_")
}
