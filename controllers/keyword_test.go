package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testFile "github.com/gutakk/go-google-scraper/tests/file"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
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

	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorHelper.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	db.SetupRedisPool()

	testDB.InitKeywordStatusEnum(db.GetDB())
	err = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if err != nil {
		log.Fatal(errorHelper.MigrateDatabaseFailure, err)
	}

	s.engine = testConfig.GetRouter(true)
	new(LoginController).applyRoutes(EnsureGuestUserGroup(s.engine))
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(s.engine))

	email := faker.Email()
	password := faker.Password()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorHelper.HashPasswordFailure, err)
	}

	user := models.User{Email: email, Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.userID = user.ID
}

func (s *KeywordDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	_, err := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "search"))
	if err != nil {
		log.Fatal(errorHelper.DeleteRedisJobFailure, err)
	}
}

func TestKeywordDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordDbTestSuite))
}

func (s *KeywordDbTestSuite) TestDisplayKeywordWithAuthenticatedUserWithoutFilter() {
	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/keyword", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, isKeywordPage)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordWithAuthenticatedUserWithFilter() {
	// Cookie from login API Set-Cookie header
	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	url := "/keyword?" +
		"filter[keyword]=Test&" +
		"filter[url]=Test" +
		"filter[is_adword_advertiser]=Test"
	response := testHttp.PerformRequest(s.engine, "GET", url, headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordPage := err == nil && strings.Index(string(p), "<title>Keyword</title>") > 0

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, isKeywordPage)
}

func TestDisplayKeywordWithGuestUser(t *testing.T) {
	engine := testConfig.GetRouter(true)
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	response := testHttp.PerformRequest(engine, "GET", "/keyword", nil, nil)

	assert.Equal(t, http.StatusFound, response.Code)
	assert.Equal(t, "/login", response.Header().Get("Location"))
}

func (s *KeywordDbTestSuite) TestDisplayKeywordWithUserIDCookieButNoUser() {
	cookie := fixture.GenerateCookie("user_id", "test-user")
	headers := http.Header{}
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/keyword", headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}

func (s *KeywordDbTestSuite) TestDisplayKeywordResultWithAuthenticatedUserAndValidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)
	keywordID := fmt.Sprint(keyword.ID)
	url := fmt.Sprintf("/keyword/%s", keywordID)

	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", url, headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isKeywordResultPage := err == nil && strings.Index(string(p), keyword.Keyword) > 0

	assert.Equal(s.T(), http.StatusOK, response.Code)
	assert.Equal(s.T(), true, isKeywordResultPage)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordResultWithAuthenticatedUserButInvalidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/keyword/invalid-keyword", headers, nil)
	p, err := ioutil.ReadAll(response.Body)
	isNotFoundPage := err == nil && strings.Index(string(p), "<title>Not Found</title>") > 0

	assert.Equal(s.T(), http.StatusNotFound, response.Code)
	assert.Equal(s.T(), true, isNotFoundPage)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordResultWithGuestUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)
	keywordID := fmt.Sprint(keyword.ID)
	url := fmt.Sprintf("/keyword/%s", keywordID)

	response := testHttp.PerformRequest(s.engine, "GET", url, nil, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}

func (s *KeywordDbTestSuite) TestDisplayKeywordHTMLWithAuthenticatedUserAndValidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name(), HtmlCode: "test-html"}
	db.GetDB().Create(&keyword)
	keywordID := fmt.Sprint(keyword.ID)
	url := fmt.Sprintf("/keyword/%s/html", keywordID)

	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", url, headers, nil)

	assert.Equal(s.T(), http.StatusOK, response.Code)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordHTMLWithAuthenticatedUserButInvalidKeyword() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name(), HtmlCode: "test-html"}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", "/keyword/invalid-keyword/html", headers, nil)

	assert.Equal(s.T(), http.StatusNotFound, response.Code)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordHTMLWithAuthenticatedUserButNoHTMLCode() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)
	keywordID := fmt.Sprint(keyword.ID)
	url := fmt.Sprintf("/keyword/%s/html", keywordID)

	headers := http.Header{}
	cookie := fixture.GenerateCookie("user_id", s.userID)
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "GET", url, headers, nil)

	assert.Equal(s.T(), http.StatusNotFound, response.Code)
}

func (s *KeywordDbTestSuite) TestDisplayKeywordHTMLWithGuestUser() {
	keyword := models.Keyword{UserID: s.userID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)
	keywordID := fmt.Sprint(keyword.ID)
	url := fmt.Sprintf("/keyword/%s/html", keywordID)

	response := testHttp.PerformRequest(s.engine, "GET", url, nil, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
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
