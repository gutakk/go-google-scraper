package controllers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/tests"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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

func TestRegisterWithValidParameters(t *testing.T) {
	engine := tests.GetRouter(true)
	gormDB, _ := gorm.Open(postgres.Open("sslmode=disable host=localhost port=5432 dbname=go_google_scraper_development user=postgres"), &gorm.Config{})
	authController := &AuthController{DB: gormDB}

	authController.applyRoutes((engine))

	formData := url.Values{}
	formData.Set("email", "test@hello.com")
	formData.Set("password", "123456")
	formData.Set("confirm-password", "123456")

	// w := tests.PerformRequest(engine, "POST", "/register")
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusFound, w.Code)
}
