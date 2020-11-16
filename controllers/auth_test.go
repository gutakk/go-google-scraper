package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/tests"
	"gopkg.in/go-playground/assert.v1"
)

func TestDisplayRegister(t *testing.T) {
	engine := tests.GetRouter(true)
	new(AuthController).applyRoutes(engine)

	w := tests.PerformRequest(engine, "GET", "/register")
	p, err := ioutil.ReadAll(w.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, pageOK)
}
