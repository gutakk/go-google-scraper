package controllers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testFile "github.com/gutakk/go-google-scraper/tests/file"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type KeywordDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
	userID uint
}

func (s *KeywordDbTestSuite) SetupTest() {
	config.LoadEnv()

	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	db.GenerateRedisPool()

	testDB.InitKeywordStatusEnum(db.GetDB())
	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})

	s.engine = testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(s.engine))
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(s.engine))

	email := faker.Email()
	password := faker.Password()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{Email: email, Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID
}

func (s *KeywordDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	_, _ = db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "scraping"))
}

func TestKeywordDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordDbTestSuite))
}

func TestDisplayKeywordWithAuthenticatedUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := "go-google-scraper=MTYwNjQ2Mjk3MXxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRFdFE9PXzl6APqAQw3gAQqlHoXMYrPpnqPFkEP8SRHJZEpl-_LDQ=="
	headers.Set("Cookie", cookie)

	response := testHttp.PerformRequest(engine, "GET", "/keyword", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, true, isKeywordPage)
}

func TestDisplayKeywordWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := testHttp.PerformRequest(engine, "GET", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func (s *KeywordDbTestSuite) TestUploadKeywordWithAuthenticatedUserAndValidParams() {
	headers, payload := testFile.CreateMultipartPayload("tests/fixture/adword_keywords.csv")
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformFileUploadRequest(s.engine, "POST", "/keyword", headers, payload)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/keyword", response.Header().Get("Location"))
}

func (s *KeywordDbTestSuite) TestUploadKeywordWithAuthenticatedUserAndBlankPayload() {
	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformFileUploadRequest(s.engine, "POST", "/keyword", headers, &bytes.Buffer{})

	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0
	pageError := err == nil && strings.Index(string(p), "File is required") > 0

	assert.Equal(s.T(), http.StatusBadRequest, response.Code)
	assert.Equal(s.T(), true, isKeywordPage)
	assert.Equal(s.T(), true, pageError)
}

func TestUploadKeywordWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := testHttp.PerformRequest(engine, "POST", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}
