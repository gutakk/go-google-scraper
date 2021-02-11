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
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
)

type RegisterDbTestSuite struct {
	suite.Suite
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
	email    string
	password string
}

func (s *RegisterDbTestSuite) SetupTest() {
	testDB.SetupTestDatabase()

	s.engine = testConfig.GetRouter(true)
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
	response := testHttp.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}

func (s *RegisterDbTestSuite) TestRegisterWithBlankEmailValidation() {
	s.formData.Del("email")

	response := testHttp.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageError := testHttp.ValidateResponseBody(bodyByte, "invalid email format")

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
}

func (s *RegisterDbTestSuite) TestRegisterWithBlankPasswordValidation() {
	s.formData.Del("password")

	response := testHttp.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageError := testHttp.ValidateResponseBody(bodyByte, "Password is required")
	isEmailFieldValueExist := testHttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestRegisterWithPasswordNotMatchValidation() {
	s.formData.Set("confirm-password", "invalid")

	response := testHttp.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageError := testHttp.ValidateResponseBody(bodyByte, "passwords do not match")
	isEmailFieldValueExist := testHttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestRegisterWithTooShortPasswordValidation() {
	s.formData.Set("password", "12345")
	s.formData.Set("confirm-password", "12345")

	response := testHttp.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageError := testHttp.ValidateResponseBody(bodyByte, "Password must be longer than 6")
	isEmailFieldValueExist := testHttp.ValidateResponseBody(bodyByte, s.email)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, pageError)
	assert.Equal(s.T(), true, isEmailFieldValueExist)
}

func (s *RegisterDbTestSuite) TestDisplayRegisterWithAuthenticatedUser() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(s.password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}
	user := models.User{Email: s.email, Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	cookie := fixture.GenerateCookie("user_id", user.ID)
	s.headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/register", s.headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func TestDisplayRegister(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))

	response := testHttp.PerformRequest(engine, "GET", "/register", nil, nil)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageOK := testHttp.ValidateResponseBody(bodyByte, "<title>Register</title>")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

func (s *RegisterDbTestSuite) TestDisplayRegisterWithUserIDCookieButNoUser() {
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers := http.Header{}
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/register", headers, nil)

	bodyByte := testHttp.ReadResponseBody(response.Body)
	pageOK := testHttp.ValidateResponseBody(bodyByte, "<title>Register</title>")

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, pageOK) // TODO: Check the controller in other task
}
