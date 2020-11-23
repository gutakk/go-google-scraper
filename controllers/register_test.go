package controllers

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

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
	new(RegisterController).applyRoutes(engine)

	response := tests.PerformRequest(engine, "GET", "/register", nil, nil)
	p, err := ioutil.ReadAll(response.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

type DBTestSuite struct {
	suite.Suite
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
}

func (s *DBTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	s.engine = tests.GetRouter(true)
	registerController := &RegisterController{}
	registerController.applyRoutes(s.engine)

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	s.formData = url.Values{}
	s.formData.Set("email", "test@hello.com")
	s.formData.Set("password", "123456")
	s.formData.Set("confirm-password", "123456")
}

func (s *DBTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}

func (s *DBTestSuite) TestRegisterWithValidParameters() {
	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	user := models.User{}
	db.GetDB().First(&user)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
	assert.Equal(s.T(), "test@hello.com", user.Email)
}

func (s *DBTestSuite) TestRegisterWithBlankEmail() {
	s.formData.Del("email")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Invalid email format") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *DBTestSuite) TestRegisterWithBlankPassword() {
	s.formData.Del("password")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password is required") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@hello.com") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *DBTestSuite) TestRegisterWithPasswordNotMatch() {
	s.formData.Set("confirm-password", "1234567")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Passwords do not match") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@hello.com") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *DBTestSuite) TestRegisterWithTooShortPassword() {
	s.formData.Set("password", "12345")
	s.formData.Set("confirm-password", "12345")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password must be longer than 6") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@hello.com") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *DBTestSuite) TestRegisterWithDuplicateEmail() {
	db.GetDB().Create(&models.User{Email: "test@hello.com", Password: "123456"})

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Email already exists") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@hello.com") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}
