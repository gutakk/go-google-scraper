package controllers

import (
	"net/http"
	"testing"

	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testdb "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testhttp "github.com/gutakk/go-google-scraper/tests/http"

	"gopkg.in/go-playground/assert.v1"
)

func TestDisplayHomeWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	response := testhttp.PerformRequest(engine, "GET", "/", nil, nil)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	isHomePage := testhttp.ValidateResponseBody(bodyByte, "<title>Home</title>")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isHomePage)
}

func TestDisplayHomeWithAuthenticatedUser(t *testing.T) {
	testdb.SetupTestDatabase()

	engine := testConfig.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testhttp.PerformRequest(engine, "GET", "/", headers, nil)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	isHomePage := testhttp.ValidateResponseBody(bodyByte, "<title>Home</title>")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isHomePage)
}
