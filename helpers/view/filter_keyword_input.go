package view

func FilterKeywordPartialInput(filter map[string][]string, label string, name string, placeholder string) map[string]interface{} {
	return map[string]interface{}{
		"filter":            filter,
		"filterLabel":       label,
		"filterName":        name,
		"filterPlaceholder": placeholder,
	}
}
