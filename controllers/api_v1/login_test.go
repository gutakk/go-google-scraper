package api_v1_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api_v1"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
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

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()
	_ = oauth.SetupOAuthServer()
	database, _ := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	db.GetDB = func() *gorm.DB {
		return database
	}

	_ = db.GetDB().AutoMigrate(&models.User{})
}

type LoginAPIControllerDbTestSuite struct {
	suite.Suite
	engine      *gin.Engine
	user        models.User
	oauthClient oauth_test.OAuthClient
	headers     http.Header
}

func (s *LoginAPIControllerDbTestSuite) SetupTest() {
	s.engine = testConfig.GetRouter(false)
	new(api_v1.LoginAPIController).ApplyRoutes(controllers.PublicAPIGroup(s.engine.Group("/api/v1")))

	s.headers = http.Header{}
	s.headers.Set("Content-Type", "application/x-www-form-urlencoded")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)
	s.user = user

	s.oauthClient = oauth_test.OAuthClient{
		ID:     "client-id",
		Secret: "client-secret",
		Domain: "http://localhost:8080",
	}
	data, _ := json.Marshal(s.oauthClient)
	s.oauthClient.Data = data

	db.GetDB().Exec("INSERT INTO oauth2_clients VALUES(?, ?, ?, ?)",
		s.oauthClient.ID,
		s.oauthClient.Secret,
		s.oauthClient.Domain,
		s.oauthClient.Data,
	)
}

func (s *LoginAPIControllerDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
	db.GetDB().Exec("DELETE FROM oauth2_clients")
	db.GetDB().Exec("DELETE FROM oauth2_tokens")
}

func TestLoginAPIControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(LoginAPIControllerDbTestSuite))
}

func (s *LoginAPIControllerDbTestSuite) TestGenerateTokenWithValidParams() {
	formData := url.Values{}
	formData.Set("username", s.user.Email)
	formData.Set("password", "password")
	formData.Set("grant_type", "password")
	formData.Set("client_id", s.oauthClient.ID)
	formData.Set("client_secret", s.oauthClient.Secret)

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/login", s.headers, formData)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string]string
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	var data []byte
	row := db.GetDB().Table("oauth2_tokens").Select("data").Row()
	_ = row.Scan(&data)

	var dataVal map[string]interface{}
	_ = json.Unmarshal(data, &dataVal)

	assert.Equal(s.T(), http.StatusOK, resp.Code)
	assert.Equal(s.T(), parsedRespBody["access_token"], dataVal["Access"])
	assert.Equal(s.T(), parsedRespBody["refresh_token"], dataVal["Refresh"])
	assert.Equal(s.T(), fmt.Sprint(s.user.ID), dataVal["UserID"])
}

func (s *LoginAPIControllerDbTestSuite) TestGenerateTokenWithInvalidGrantType() {
	formData := url.Values{}
	formData.Set("username", s.user.Email)
	formData.Set("password", "password")
	formData.Set("grant_type", "invalidGrant")
	formData.Set("client_id", s.oauthClient.ID)
	formData.Set("client_secret", s.oauthClient.Secret)

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/login", s.headers, formData)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string]string
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrUnsupportedGrantType.Error(), parsedRespBody["error"])
}

func (s *LoginAPIControllerDbTestSuite) TestGenerateTokenWithInvalidClientID() {
	formData := url.Values{}
	formData.Set("username", s.user.Email)
	formData.Set("password", "password")
	formData.Set("grant_type", "password")
	formData.Set("client_id", "invalid")
	formData.Set("client_secret", s.oauthClient.Secret)

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/login", s.headers, formData)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string]string
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusInternalServerError, resp.Code)
	assert.Equal(s.T(), errors.ErrServerError.Error(), parsedRespBody["error"])
}

func (s *LoginAPIControllerDbTestSuite) TestGenerateTokenWithInvalidClientSecret() {
	formData := url.Values{}
	formData.Set("username", s.user.Email)
	formData.Set("password", "password")
	formData.Set("grant_type", "password")
	formData.Set("client_id", s.oauthClient.ID)
	formData.Set("client_secret", "invalid")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/login", s.headers, formData)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string]string
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
	assert.Equal(s.T(), errors.ErrInvalidClient.Error(), parsedRespBody["error"])
}