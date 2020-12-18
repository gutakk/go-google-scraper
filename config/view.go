package config

import (
	"html/template"

	"github.com/foolin/goview"
	"github.com/gutakk/go-google-scraper/helpers/view"
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
			"scrapeResultPartial": view.ScrapeResultPartial,
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
