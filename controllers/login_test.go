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
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type LoginDbTestSuite struct {
	suite.Suite
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
	email    string
	password string
}

func (s *LoginDbTestSuite) SetupTest() {
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	_ = db.GetDB().AutoMigrate(&models.User{})

	s.engine = testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(s.engine))

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

func (s *LoginDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestLoginDbTestSuite(t *testing.T) {
	suite.Run(t, new(LoginDbTestSuite))
}

func (s *LoginDbTestSuite) TestLoginWithValidParameters() {
	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *LoginDbTestSuite) TestDisplayLoginWithAuthenticatedUser() {
	cookie := fixture.GenerateCookie("user_id", "test-user")
	s.headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/login", s.headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *LoginDbTestSuite) TestLoginWithBlankEmailValidation() {
	s.formData.Del("email")

	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "invalid email format") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *LoginDbTestSuite) TestLoginWithBlankPasswordValidation() {
	s.formData.Del("password")

	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password is required") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithTooShortPasswordValidation() {
	s.formData.Set("password", "12345")

	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "Password must be longer than 6") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithInvalidEmail() {
	s.formData.Set("email", "test@email.com")

	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "username or password is invalid") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), "test@email.com") > 0

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithInvalidPassword() {
	s.formData.Set("password", "123456789")

	response := testHttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)
	p, err := ioutil.ReadAll(response.Body)
	pageError := err == nil && strings.Index(string(p), "username or password is invalid") > 0
	isEmailFieldValueExist := err == nil && strings.Index(string(p), s.email) > 0

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func TestDisplayLogin(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))

	response := testHttp.PerformRequest(engine, "GET", "/login", nil, nil)
	p, err := ioutil.ReadAll(response.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Login</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

func (s *LoginDbTestSuite) TestDisplayLoginWithUserIDCookieButNoUser() {
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers := http.Header{}
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/login", headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}
