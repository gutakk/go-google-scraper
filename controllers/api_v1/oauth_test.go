package api_v1_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api_v1"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
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
}

type OAuthControllerDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
}

func (s *OAuthControllerDbTestSuite) SetupTest() {
	s.engine = testConfig.GetRouter(true)
	new(api_v1.OAuthController).ApplyRoutes(controllers.BasicAuthAPIGroup(s.engine.Group("/api")))
}

func (s *OAuthControllerDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM oauth2_clients")
}

func TestOAuthControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthControllerDbTestSuite))
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithValidBasicAuth() {
	headers := http.Header{}
	// Basic auth with username = admin and password = password
	headers.Set("Authorization", "Basic YWRtaW46cGFzc3dvcmQ=")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/client", headers, nil)
	respBodyData, _ := ioutil.ReadAll(resp.Body)
	var parsedRespBody map[string]api_helper.DataResponseObject
	_ = json.Unmarshal(respBodyData, &parsedRespBody)

	var data []byte
	row := db.GetDB().Table("oauth2_clients").Select("data").Row()
	_ = row.Scan(&data)

	var dataVal map[string]interface{}
	_ = json.Unmarshal(data, &dataVal)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
	assert.Equal(s.T(), parsedRespBody["data"].Attributes["client_id"], dataVal["ID"])
	assert.Equal(s.T(), parsedRespBody["data"].Attributes["client_secret"], dataVal["Secret"])
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithInvalidBasicAuth() {
	headers := http.Header{}
	// Basic auth with username = admin and password = password
	headers.Set("Authorization", "Basic invalid")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/client", headers, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithoutBasicAuth() {
	resp := testHttp.PerformRequest(s.engine, "POST", "/api/client", nil, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}
