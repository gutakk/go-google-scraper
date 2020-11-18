package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
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

	w := tests.PerformRequest(engine, "GET", "/register")
	p, err := ioutil.ReadAll(w.Body)
	pageOK := err == nil && strings.Index(string(p), "<title>Register</title>") > 0

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, true, pageOK)
}

type DBTestSuite struct {
	suite.Suite
	DB       *gorm.DB
	engine   *gin.Engine
	formData url.Values
}

func (s *DBTestSuite) SetupTest() {
	db, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	s.DB = db

	if error := db.AutoMigrate(&models.User{}); error != nil {
		log.Fatal(fmt.Sprintf("Failed to migrate database %v", error))
	} else {
		log.Print("Migrate to database successfully")
	}

	s.engine = tests.GetRouter(true)
	authController := &AuthController{DB: s.DB}
	authController.applyRoutes(s.engine)

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

	w := httptest.NewRecorder()
	s.engine.ServeHTTP(w, req)

	user := models.User{}
	s.DB.First(&user)

	assert.Equal(s.T(), http.StatusFound, w.Code)
	assert.Equal(s.T(), "/", w.Header().Get("Location"))
	assert.Equal(s.T(), "test@hello.com", user.Email)
}

func TestDBTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
