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
	DB *gorm.DB
}

func (suite *DBTestSuite) SetupTest() {
	db, _ := gorm.Open(postgres.Open(tests.ConstructTestDsn()), &gorm.Config{})
	suite.DB = db

	if error := db.AutoMigrate(&models.User{}); error != nil {
		log.Fatal(fmt.Sprintf("Failed to migrate database %v", error))
	} else {
		log.Print("Migrate to database successfully")
	}
}

func (suite *DBTestSuite) TearDownTest() {
	suite.DB.Exec("DELETE FROM users")
}

func (suite *DBTestSuite) TestRegisterWithValidParameters() {
	engine := tests.GetRouter(true)
	authController := &AuthController{DB: suite.DB}

	authController.applyRoutes((engine))

	formData := url.Values{}
	formData.Set("email", "test@hello.com")
	formData.Set("password", "123456")
	formData.Set("confirm-password", "123456")

	req, _ := http.NewRequest("POST", "/register", strings.NewReader(formData.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	w := httptest.NewRecorder()
	engine.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusFound, w.Code)
}

func TestTestSuite(t *testing.T) {
	suite.Run(t, new(DBTestSuite))
}
