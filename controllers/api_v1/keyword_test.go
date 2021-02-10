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
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api_v1"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testFile "github.com/gutakk/go-google-scraper/tests/file"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/oauth_test"
	"github.com/gutakk/go-google-scraper/tests/path_test"

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
		log.Fatal(errorconf.ChangeToRootDirFailure, err)
	}

	config.LoadEnv()

	err = oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(errorconf.StartOAuthServerFailute, err)
	}

	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(errorconf.ConnectToDatabaseFailure, err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	testDB.InitKeywordStatusEnum(db.GetDB())
	err = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
	if err != nil {
		log.Fatal(errorconf.MigrateDatabaseFailure, err)
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
		log.Error(errorconf.HashPasswordFailure, err)
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
		log.Error(errorconf.JSONMarshalFailure, err)
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
		log.Fatal(errorconf.DeleteRedisJobFailure, err)
	}
}

func TestKeywordAPIControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordAPIControllerDbTestSuite))
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordWithValidParams() {
	keyword := models.Keyword{
		Model:                   &gorm.Model{ID: 1},
		Keyword:                 "testKeyword",
		Status:                  models.Pending,
		LinksCount:              1,
		NonAdwordsCount:         1,
		TopPositionAdwordsCount: 1,
		TotalAdwordsCount:       1,
		UserID:                  s.user.ID,
		HtmlCode:                "testHTML",
		FailedReason:            "",
	}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/1", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string]api_helper.DataResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	data := parsedRespBody["data"]
	attributes, _ := parsedRespBody["data"].Attributes.(map[string]interface{})
	relationships := parsedRespBody["data"].Relationships

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), data.ID, "1")
	assert.Equal(s.T(), data.Type, "keyword")
	assert.Equal(s.T(), attributes["keyword"], "testKeyword")
	assert.Equal(s.T(), attributes["status"], "pending")
	assert.Equal(s.T(), attributes["links_count"], float64(1))
	assert.Equal(s.T(), attributes["non_adwords_count"], float64(1))
	assert.Equal(s.T(), attributes["non_adword_links"], nil)
	assert.Equal(s.T(), attributes["top_position_adwords_count"], float64(1))
	assert.Equal(s.T(), attributes["top_position_adword_links"], nil)
	assert.Equal(s.T(), attributes["total_adwords_count"], float64(1))
	assert.Equal(s.T(), attributes["html_code"], "testHTML")
	assert.Equal(s.T(), attributes["failed_reason"], "")
	assert.Equal(s.T(), relationships["user"].Data.Type, "user")
	assert.Equal(s.T(), relationships["user"].Data.ID, fmt.Sprint(s.user.ID))
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordWithInvalidKeywordID() {
	keyword := models.Keyword{
		Model:                   &gorm.Model{ID: 1},
		Keyword:                 "testKeyword",
		Status:                  models.Pending,
		LinksCount:              1,
		NonAdwordsCount:         1,
		TopPositionAdwordsCount: 1,
		TotalAdwordsCount:       1,
		UserID:                  s.user.ID,
		HtmlCode:                "testHTML",
		FailedReason:            "",
	}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/9999", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "keyword not found", parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordWithValidParamsButNotTheResourceOwner() {
	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-not-resource-owner",
		Refresh:   "test-refresh",
	}

	data, err := json.Marshal(&oauth_test.TokenData{
		Access: tokenItem.Access,
		UserID: "invalidUserID",
	})
	if err != nil {
		log.Error(errorconf.JSONMarshalFailure, err)
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

	keyword := models.Keyword{
		Model:                   &gorm.Model{ID: 1},
		Keyword:                 "testKeyword",
		Status:                  models.Pending,
		LinksCount:              1,
		NonAdwordsCount:         1,
		TopPositionAdwordsCount: 1,
		TotalAdwordsCount:       1,
		UserID:                  s.user.ID,
		HtmlCode:                "testHTML",
		FailedReason:            "",
	}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-not-resource-owner")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/1", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusNotFound, resp.Code)
	assert.Equal(s.T(), "keyword not found", parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordAPIWithoutAuthorizationHeader() {
	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/1", nil, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordAPIWithInvalidAccessToken() {
	headers := http.Header{}
	headers.Set("Authorization", "invalid_token")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/1", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordAPIWithExpiredAccessToken() {
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
		log.Error(errorconf.JSONMarshalFailure, err)
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

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords/1", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrExpiredAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsWithValidParams() {
	keyword := models.Keyword{UserID: s.user.ID, Keyword: "testKeyword"}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.DataResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	respValue, _ := parsedRespBody["data"][0].Attributes.(map[string]interface{})

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), 1, len(parsedRespBody["data"]))
	assert.Equal(s.T(), "testKeyword", respValue["keyword"])
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsWithValidParamsButNoKeywords() {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.DataResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), 0, len(parsedRespBody["data"]))
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsWithValidParamsButNotTheResourceOwner() {
	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-not-resource-owner",
		Refresh:   "test-refresh",
	}

	data, err := json.Marshal(&oauth_test.TokenData{
		Access: tokenItem.Access,
		UserID: "invalidUserID",
	})
	if err != nil {
		log.Error(errorconf.JSONMarshalFailure, err)
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

	keyword := models.Keyword{UserID: s.user.ID, Keyword: faker.Name()}
	db.GetDB().Create(&keyword)

	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-not-resource-owner")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.DataResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), 0, len(parsedRespBody["data"]))
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsAPIWithoutAuthorizationHeader() {
	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", nil, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsAPIWithInvalidAccessToken() {
	headers := http.Header{}
	headers.Set("Authorization", "invalid_token")

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestFetchKeywordsAPIWithExpiredAccessToken() {
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
		log.Error(errorconf.JSONMarshalFailure, err)
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

	resp := testHttp.PerformRequest(s.engine, "GET", "/api/v1/keywords", headers, nil)
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrExpiredAccessToken.Error(), parsedRespBody["errors"][0].Detail)
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
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
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
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
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
		log.Error(errorconf.JSONMarshalFailure, err)
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
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
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
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
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
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), "invalid file", parsedRespBody["errors"][0].Detail)
}
