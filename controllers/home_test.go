package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/tests"
	"gopkg.in/go-playground/assert.v1"
)

func TestDisplayHome(t *testing.T) {
	engine := tests.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	w := tests.PerformRequest(engine, "GET", "/", nil, nil)
	p, err := ioutil.ReadAll(w.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Home</title>") > 0

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, pageOK)
}
