package api_v1_test

import (
	"bytes"
	"encoding/json"
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

	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()
	_ = oauth.SetupOAuthServer()
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	_ = db.GetDB().AutoMigrate(&models.User{}, &models.Keyword{})
}

type KeywordAPIControllerDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
	// user        models.User
	// oauthClient oauth_test.OAuthClient
}

func (s *KeywordAPIControllerDbTestSuite) SetupTest() {
	s.engine = testConfig.GetRouter(false)
	new(api_v1.KeywordAPIController).ApplyRoutes(controllers.PrivateAPIGroup(s.engine.Group("/api/v1")))

	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-access",
		Refresh:   "test-refresh",
	}

	data, _ := json.Marshal(tokenItem)
	tokenItem.Data = data

	db.GetDB().Exec("INSERT INTO oauth2_tokens(created_at, expires_at, code, access, refresh, data) VALUES(?, ?, ?, ?, ?, ?)",
		tokenItem.CreatedAt,
		tokenItem.ExpiresAt,
		tokenItem.Code,
		tokenItem.Access,
		tokenItem.Refresh,
		tokenItem.Data,
	)

	// hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	// user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	// db.GetDB().Create(&user)
	// s.user = user

	// s.oauthClient = oauth_test.OAuthClient{
	// 	ID:     "client-id",
	// 	Secret: "client-secret",
	// 	Domain: "http://localhost:8080",
	// }
	// data, _ := json.Marshal(s.oauthClient)
	// s.oauthClient.Data = data

	// db.GetDB().Exec("INSERT INTO oauth2_clients VALUES(?, ?, ?, ?)",
	// 	s.oauthClient.ID,
	// 	s.oauthClient.Secret,
	// 	s.oauthClient.Domain,
	// 	s.oauthClient.Data,
	// )
}

func (s *KeywordAPIControllerDbTestSuite) TearDownTest() {
	// db.GetDB().Exec("DELETE FROM users")
	db.GetDB().Exec("DELETE FROM oauth2_tokens")
}

func TestKeywordAPIControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(KeywordAPIControllerDbTestSuite))
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordWithoutAuthorizationHeader() {
	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keyword", nil, nil)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordWithInvalidAccessToken() {
	headers := http.Header{}
	headers.Set("Authorization", "invalid_token")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keyword", headers, nil)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordWithExpiredAccessToken() {
	db.GetDB().Exec("DELETE FROM oauth2_tokens")

	tokenItem := &oauth_test.TokenStoreItem{
		CreatedAt: time.Now(),
		ExpiresAt: time.Now(),
		Code:      "test-code",
		Access:    "test-access",
		Refresh:   "test-refresh",
	}

	data, _ := json.Marshal(&oauth_test.TokenData{
		AccessExpiresIn: 1,
		Access:          tokenItem.Access,
	})
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

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keyword", headers, nil)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrExpiredAccessToken.Error(), parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordWithAccessTokenButNoFileFormInTheRequest() {
	headers := http.Header{}
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/keyword", headers, nil)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), "invalid file", parsedRespBody["errors"][0].Detail)
}

func (s *KeywordAPIControllerDbTestSuite) TestUploadKeywordWithAccessTokenButBlankPayload() {
	headers, _ := testFile.CreateMultipartPayload("")
	headers.Set("Authorization", "Bearer test-access")

	resp := testHttp.PerformFileUploadRequest(s.engine, "POST", "/api/v1/keyword", headers, &bytes.Buffer{})
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string][]api_helper.ErrorResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusBadRequest, resp.Code)
	assert.Equal(s.T(), "invalid file", parsedRespBody["errors"][0].Detail)
}
