package view

func FilterKeywordPartial(query map[string][]string, label string, name string, placeholder string) map[string]interface{} {
	return map[string]interface{}{
		"query":             query,
		"filterLabel":       label,
		"filterName":        name,
		"filterPlaceholder": placeholder,
	}
}
