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
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDisplayLogin(t *testing.T) {
	engine := tests.GetRouter(true)
	new(UserSessionController).applyRoutes(EnsureNoAuthenticationGroup(engine))

	response := tests.PerformRequest(engine, "GET", "/login", nil, nil)
	p, err := ioutil.ReadAll(response.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Login</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

type UserSessionDbTestSuite struct {
	suite.Suite
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
	email    string
	password string
}

func (s *UserSessionDbTestSuite) SetupTest() {
	testDB, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return testDB
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	s.engine = tests.GetRouter(true)
	new(UserSessionController).applyRoutes(EnsureNoAuthenticationGroup(s.engine))

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	s.email = faker.Email()
	s.password = faker.Password()

	s.formData = url.Values{}
	s.formData.Set("email", s.email)
	s.formData.Set("password", s.password)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
	db.GetDB().Create(&models.User{Email: s.email, Password: string(hashedPassword)})
}

func (s *UserSessionDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestUserSessionDbTestSuite(t *testing.T) {
	suite.Run(t, new(UserSessionDbTestSuite))
}

func (s *UserSessionDbTestSuite) TestLoginWithValidParameters() {
	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *UserSessionDbTestSuite) TestDisplayLoginWithUserSession() {
	// Cookie from login API Set-Cookie header
	cookie := "mysession=MTYwNjI3ODk0NHxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRCR0E9PXxa_dKXde8j6m4z_kPgaiPYuDGHj79HxhCMNw3zIoeM6g=="
	s.headers.Set("Cookie", cookie)

	response := tests.PerformRequest(s.engine, "GET", "/login", s.headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *UserSessionDbTestSuite) TestLoginWithBlankEmailValidation() {
	s.formData.Del("email")

	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Invalid email format") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *UserSessionDbTestSuite) TestLoginWithBlankPasswordValidation() {
	s.formData.Del("password")

	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password is required") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *UserSessionDbTestSuite) TestLoginWithTooShortPasswordValidation() {
	s.formData.Set("password", "12345")

	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password must be longer than 6") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *UserSessionDbTestSuite) TestLoginWithInvalidEmail() {
	s.formData.Set("email", "test@email.com")

	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Username or password is invalid") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@email.com") > 0

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *UserSessionDbTestSuite) TestLoginWithInvalidPassword() {
	s.formData.Set("password", "123456789")

	response := tests.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Username or password is invalid") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}
