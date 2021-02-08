package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	log "github.com/sirupsen/logrus"
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
	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorHelper.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	err = db.GetDB().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(errorHelper.MigrateDatabaseFailure, err)
	}

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
