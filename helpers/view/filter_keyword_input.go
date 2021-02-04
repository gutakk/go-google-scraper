package view

func FilterKeywordPartialInput(filter map[string][]string, label string, inputType string, className string, name string, placeholder string) map[string]interface{} {
	var filterValue interface{}
	if len(filter) > 0 {
		switch inputType {
		case "text":
			filterValue = filter[name][0]
		case "checkbox":
			filterValue = filter[name]
		}
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
