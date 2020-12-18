package config

import (
	"html/template"

	"github.com/foolin/goview"
)

func AppGoviewConfig() goview.Config {
	return goview.Config{
		Root:         "templates",
		Extension:    ".html",
		Master:       "layouts/application",
		DisableCache: false,
		Partials: []string{
			"partials/scrape_result",
			"partials/scrape_result_list",
		},
		Funcs: template.FuncMap{
			"scrapeResult": func(title string, value interface{}) map[string]interface{} {
				return map[string]interface{}{
					"title": title,
					"value": value,
				}
			},
		},
	}
}

func AuthenticationGoviewConfig() goview.Config {
	return goview.Config{
		Root:         "templates",
		Extension:    ".html",
		Master:       "layouts/authentication",
		DisableCache: false,
	}
}
