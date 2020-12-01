package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/tests"
	"gopkg.in/go-playground/assert.v1"
)

func TestDisplayKeywordWithGuestUser(t *testing.T) {
	engine := tests.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := tests.PerformRequest(engine, "GET", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func TestDisplayKeywordWithAuthenticatedUser(t *testing.T) {
	engine := tests.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := "go-google-scraper=MTYwNjQ2Mjk3MXxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRFdFE9PXzl6APqAQw3gAQqlHoXMYrPpnqPFkEP8SRHJZEpl-_LDQ=="
	headers.Set("Cookie", cookie)

	response := tests.PerformRequest(engine, "GET", "/keyword", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isKeywordPage)
}
