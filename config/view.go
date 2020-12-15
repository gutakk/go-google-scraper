package config

import (
	"html/template"

	"github.com/foolin/goview"
	"github.com/gutakk/go-google-scraper/presenters"
)

func AppGoviewConfig() goview.Config {
	return goview.Config{
		Root:         "templates",
		Extension:    ".html",
		Master:       "layouts/application",
		DisableCache: false,
		Funcs: template.FuncMap{
			"formattedCreatedAt": presenters.FormattedCreatedAt,
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
