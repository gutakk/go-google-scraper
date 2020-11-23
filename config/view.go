package config

import (
	"html/template"

	"github.com/foolin/goview"
)

func GoviewConfig() goview.Config {
	return goview.Config{
		Root:      "templates",
		Extension: ".html",
		Master:    "layouts/master",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
		},
		DisableCache: false,
	}
}
