package view

func FilterKeywordPartialInput(filter map[string][]string, label string, inputType string, className string, name string, placeholder string) map[string]interface{} {
	var filterValue string
	if len(filter) > 0 {
		filterValue = filter[name][0]
	}

	return map[string]interface{}{
		"filter":            filterValue,
		"filterLabel":       label,
		"inputType":         inputType,
		"className":         className,
		"filterName":        name,
		"filterPlaceholder": placeholder,
	}
}
