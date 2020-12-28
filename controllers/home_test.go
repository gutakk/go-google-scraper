package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDisplayHomeWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	w := testHttp.PerformRequest(engine, "GET", "/", nil, nil)
	p, err := ioutil.ReadAll(w.Body)
	isHomePage := err == nil && strings.Index(string(p), "<title>Home</title>") > 0

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, isHomePage)
}

func TestDisplayHomeWithAuthenticatedUser(t *testing.T) {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	engine := testConfig.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(engine, "GET", "/", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isHomePage := err == nil && strings.Index(string(p), "<title>Home</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isHomePage)
}
