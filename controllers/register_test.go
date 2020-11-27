package controllers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDisplayRegister(t *testing.T) {
	engine := tests.GetRouter(true)
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))

	response := tests.PerformRequest(engine, "GET", "/register", nil, nil)
	p, err := ioutil.ReadAll(response.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

type RegisterDbTestSuite struct {
	suite.Suite
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
	email    string
	password string
}

func (s *RegisterDbTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	s.engine = tests.GetRouter(true)
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(s.engine))

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	s.email = faker.Email()
	s.password = faker.Password()

	s.formData = url.Values{}
	s.formData.Set("email", s.email)
	s.formData.Set("password", s.password)
	s.formData.Set("confirm-password", s.password)
}

func (s *RegisterDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestRegisterDbTestSuite(t *testing.T) {
	suite.Run(t, new(RegisterDbTestSuite))
}

func (s *RegisterDbTestSuite) TestRegisterWithValidParameters() {
	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}

func (s *RegisterDbTestSuite) TestRegisterWithBlankEmailValidation() {
	s.formData.Del("email")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Invalid email format") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *RegisterDbTestSuite) TestRegisterWithBlankPasswordValidation() {
	s.formData.Del("password")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password is required") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestRegisterWithPasswordNotMatchValidation() {
	s.formData.Set("confirm-password", "invalid")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Passwords do not match") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestRegisterWithTooShortPasswordValidation() {
	s.formData.Set("password", "12345")
	s.formData.Set("confirm-password", "12345")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password must be longer than 6") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestDisplayRegisterWithAuthenticatedUser() {
	// Cookie from login API Set-Cookie header
	cookie := "mysession=MTYwNjI3ODk0NHxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRCR0E9PXxa_dKXde8j6m4z_kPgaiPYuDGHj79HxhCMNw3zIoeM6g=="
	s.headers.Set("Cookie", cookie)

	response := tests.PerformRequest(s.engine, "GET", "/register", s.headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}
