package config

import (
	"html/template"

	"github.com/foolin/goview"
)

func AppGoviewConfig() goview.Config {
	return goview.Config{
		Root:      "templates",
		Extension: ".html",
		Master:    "layouts/application",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
		},
		DisableCache: false,
	}
}

func AuthenticationGoviewConfig() goview.Config {
	return goview.Config{
		Root:      "templates",
		Extension: ".html",
		Master:    "layouts/authentication",
		Partials:  []string{},
		Funcs: template.FuncMap{
			"sub": func(a, b int) int {
				return a - b
			},
		},
		DisableCache: false,
	}
}
