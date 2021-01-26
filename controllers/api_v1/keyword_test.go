package api_v1_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api_v1"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testFile "github.com/gutakk/go-google-scraper/tests/file"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/oauth_test"
	"github.com/gutakk/go-google-scraper/tests/path_test"
	log "github.com/sirupsen/logrus"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		log.Fatal(err)
	}

	config.LoadEnv()
	err = oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(err)
	}

	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	testDB.InitKeywordStatusEnum(db.GetDB())
	err = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if err != nil {
		log.Fatal(err)
	}
}

type KeywordAPIControllerDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
	user   models.User
}

func (s *KeywordAPIControllerDbTestSuite) SetupTest() {
	db.SetupRedisPool()

	s.engine = testConfig.GetRouter(false)
	new(api_v1.KeywordAPIController).ApplyRoutes(controllers.PrivateAPIGroup(s.engine.Group("/api/v1")))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
	}

	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.user = user

	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-access",
		Refresh:   "test-refresh",
	}

	data, err := json.Marshal(&oauth_test.TokenData{
		Access: tokenItem.Access,
		UserID: fmt.Sprint(s.user.ID),
	})
	if err != nil {
		log.Error(err)
	}
	tokenItem.Data = data

	db.GetDB().Exec("INSERT INTO oauth2_tokens(created_at, expires_at, code, access, refresh, data) VALUES(?, ?, ?, ?, ?, ?)",
		tokenItem.CreatedAt,
		tokenItem.ExpiresAt,
		tokenItem.Code,
		tokenItem.Access,
		tokenItem.Refresh,
		tokenItem.Data,
	)
}

func (s *KeywordAPIControllerDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM keywords")
	db.GetDB().Exec("DELETE FROM users")
	db.GetDB().Exec("DELETE FROM oauth2_tokens")
	_, err := db.GetRedisPool().Get().Do("DEL", testDB.RedisKeyJobs("go-google-scraper", "search"))
	if err != nil {
		log.Fatal(err)
	}
}

func TestKeywordAPIControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordAPIControllerDbTestSuite))
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithValidParams() {
	headers, payload := testFile.CreateMultipartPayload("tests/fixture/adword_keywords.csv")
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformFileUploadRequest(s.engine, "POST", "/api/v1/keywords", headers, payload)

	assert.Equal(s.T(), http.StatusNoContent, resp.Code)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithoutAuthorizationHeader() {
	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keywords", nil, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithInvalidAccessToken() {
	headers := http.Header{}
	headers.Set("Authorization", "invalid_token")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithExpiredAccessToken() {
	db.GetDB().Exec("DELETE FROM oauth2_tokens")

	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-access",
		Refresh:   "test-refresh",
	}

	data, err := json.Marshal(&oauth_test.TokenData{
		AccessExpiresIn: 1,
		Access:          tokenItem.Access,
	})
	if err != nil {
		log.Error(err)
	}
	tokenItem.Data = data

	db.GetDB().Exec("INSERT INTO oauth2_tokens(created_at, expires_at, code, access, refresh, data) VALUES(?, ?, ?, ?, ?, ?)",
		tokenItem.CreatedAt,
		tokenItem.ExpiresAt,
		tokenItem.Code,
		tokenItem.Access,
		tokenItem.Refresh,
		tokenItem.Data,
	)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrExpiredAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithAccessTokenButNoFileFormInTheRequest() {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), "invalid file", parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordAPIWithAccessTokenButBlankPayload() {
	headers, _ := testFile.CreateMultipartPayload("")
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformFileUploadRequest(s.engine, "POST", "/api/v1/keywords", headers, &bytes.Buffer{})
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), "invalid file", parsedRespBody["errors"][0].Detail)
}
