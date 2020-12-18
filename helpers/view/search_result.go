package view

func SearchResultPartial(title string, value interface{}) map[string]interface{} {
	return map[string]interface{}{
		"title": title,
		"value": value,
	}
}
