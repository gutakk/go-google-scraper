package controllers

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/tests"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestDisplayRegister(t *testing.T) {
	engine := tests.GetRouter(true)
	new(AuthController).applyRoutes(engine)

	response := tests.PerformRequest(engine, "GET", "/register", nil, nil)
	p, err := ioutil.ReadAll(response.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, pageOK)
}

type DBTestSuite struct {
	suite.Suite
	DB       *gorm.DB
	engine   *gin.Engine
	formData url.Values
	headers  http.Header
}

func (s *DBTestSuite) SetupTest() {
	db, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	s.DB = db

	_ = db.AutoMigrate(&models.User{})

	s.engine = tests.GetRouter(true)
	authController := &AuthController{DB: s.DB}
	authController.applyRoutes(s.engine)

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	s.formData = url.Values{}
	s.formData.Set("email", "test@hello.com")
	s.formData.Set("password", "123456")
	s.formData.Set("confirm-password", "123456")
}

func (s *DBTestSuite) TearDownTest() {
	s.DB.Exec("DELETE FROM users")
}

func (s *DBTestSuite) TestRegisterWithValidParameters() {
	req, _ := http.NewRequest("POST", "/register", strings.NewReader(s.formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	user := models.User{}
	s.DB.First(&user)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
	assert.Equal(s.T(), "test@hello.com", user.Email)
}

func (s *DBTestSuite) TestRegisterWithBlankEmail() {
	s.formData.Del("email")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)

	user := models.User{}
	result := s.DB.First(&user)

	assert.Equal(s.T(), true, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func (s *DBTestSuite) TestRegisterWithBlankPassword() {
	s.formData.Del("password")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)

	user := models.User{}
	result := s.DB.First(&user)

	assert.Equal(s.T(), true, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func (s *DBTestSuite) TestRegisterWithBlankConfirmPassword() {
	s.formData.Del("confirm-password")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)

	user := models.User{}
	result := s.DB.First(&user)

	assert.Equal(s.T(), true, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func (s *DBTestSuite) TestRegisterWithPasswordNotMatch() {
	s.formData.Set("confirm-password", "1234567")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)

	user := models.User{}
	result := s.DB.First(&user)

	assert.Equal(s.T(), true, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func (s *DBTestSuite) TestRegisterWithPasswordNotReachMinLength() {
	s.formData.Set("password", "12345")
	s.formData.Set("confirm-password", "12345")

	response := tests.PerformRequest(s.engine, "POST", "/register", s.headers, s.formData)

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)

	user := models.User{}
	result := s.DB.First(&user)

	assert.Equal(s.T(), true, errors.Is(result.Error, gorm.ErrRecordNotFound))
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
