package config

import (
	"html/template"

	"github.com/gutakk/go-google-scraper/helpers/view"

	"github.com/foolin/goview"
)

func AppGoviewConfig() goview.Config {
	return goview.Config{
		Root:         "templates",
		Extension:    ".html",
		Master:       "layouts/application",
		DisableCache: false,
		Partials: []string{
			"partials/search_result",
			"partials/list_search_result",
			"partials/filter_keyword_input",
		},
		Funcs: template.FuncMap{
			"searchResultPartial":       view.SearchResultPartial,
			"filterKeywordPartialInput": view.FilterKeywordPartialInput,
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

func ErrorGoviewConfig() goview.Config {
	return goview.Config{
		Root:         "templates",
		Extension:    ".html",
		Master:       "layouts/error",
		DisableCache: false,
	}
}
