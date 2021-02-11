package controllers

import (
	"net/http"
	"net/url"
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testhttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
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

	s.engine = testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(s.engine))

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	s.email = faker.Email()
	s.password = faker.Password()

	s.formData = url.Values{}
	s.formData.Set("email", s.email)
	s.formData.Set("password", s.password)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	db.GetDB().Create(&models.User{Email: s.email, Password: string(hashedPassword)})
}

func (s *LoginDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
}

func TestLoginDbTestSuite(t *testing.T) {
	suite.Run(t, new(LoginDbTestSuite))
}

func (s *LoginDbTestSuite) TestLoginWithValidParameters() {
	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *LoginDbTestSuite) TestDisplayLoginWithAuthenticatedUser() {
	user := models.User{}
	_ = db.GetDB().First(&user)

	cookie := fixture.GenerateCookie("user_id", user.ID)
	s.headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testhttp.PerformRequest(s.engine, "GET", "/login", s.headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *LoginDbTestSuite) TestLoginWithBlankEmailValidation() {
	s.formData.Del("email")

	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageError := testhttp.ValidateResponseBody(bodyByte, "invalid email format")

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *LoginDbTestSuite) TestLoginWithBlankPasswordValidation() {
	s.formData.Del("password")

	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageError := testhttp.ValidateResponseBody(bodyByte, "Password is required")
	isEmailFieldValueExist := testhttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithTooShortPasswordValidation() {
	s.formData.Set("password", "12345")

	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageError := testhttp.ValidateResponseBody(bodyByte, "Password must be longer than 6")
	isEmailFieldValueExist := testhttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithInvalidEmail() {
	s.formData.Set("email", "test@email.com")

	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageError := testhttp.ValidateResponseBody(bodyByte, "username or password is invalid")
	isEmailFieldValueExist := testhttp.ValidateResponseBody(bodyByte, "test@email.com")

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *LoginDbTestSuite) TestLoginWithInvalidPassword() {
	s.formData.Set("password", "123456789")

	response := testhttp.PerformRequest(s.engine, "POST", "/login", s.headers, s.formData)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageError := testhttp.ValidateResponseBody(bodyByte, "username or password is invalid")
	isEmailFieldValueExist := testhttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusUnauthorized, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func TestDisplayLogin(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))

	response := testhttp.PerformRequest(engine, "GET", "/login", nil, nil)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageOK := testhttp.ValidateResponseBody(bodyByte, "<title>Login</title>")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

func (s *LoginDbTestSuite) TestDisplayLoginWithUserIDCookieButNoUser() {
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers := http.Header{}
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testhttp.PerformRequest(s.engine, "GET", "/login", headers, nil)

	bodyByte := testhttp.ReadResponseBody(response.Body)
	pageOK := testhttp.ValidateResponseBody(bodyByte, "<title>Login</title>")

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, pageOK) // TODO: Check the controller in other task
}
