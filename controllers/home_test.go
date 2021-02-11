package controllers

import (
	"net/http"
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testhttp "github.com/gutakk/go-google-scraper/tests/http"

	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorconf.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	err = db.GetDB().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(errorconf.MigrateDatabaseFailure, err)
	}

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
