package controllers

import (
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/tests"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDisplayHomeWithoutUserSession(t *testing.T) {
	engine := tests.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	w := tests.PerformRequest(engine, "GET", "/", nil, nil)
	p, err := ioutil.ReadAll(w.Body)
	isHomePage := err == nil && strings.Index(string(p), "<title>Home</title>") > 0

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, isHomePage)
}

func TestDisplayHomeWithUserSession(t *testing.T) {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	engine := tests.GetRouter(true)
	new(HomeController).applyRoutes(engine)

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := "mysession=MTYwNjI3ODk0NHxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRCR0E9PXxa_dKXde8j6m4z_kPgaiPYuDGHj79HxhCMNw3zIoeM6g=="
	headers.Set("Cookie", cookie)

	response := tests.PerformRequest(engine, "GET", "/", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isHomePage := err == nil && strings.Index(string(p), "<title>Home</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isHomePage)
}
