package api_v1_test

import (
	"net/http"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	testJson "github.com/gutakk/go-google-scraper/tests/json"
	testPath "github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.TestMode)

	testPath.ChangeToRootDir()

	config.LoadEnv()

	testConfig.SetupTestOAuthServer()

	testDB.SetupTestDatabase()
}

type OAuthControllerDbTestSuite struct {
	suite.Suite
	engine *gin.Engine
}

func (s *OAuthControllerDbTestSuite) SetupTest() {
	s.engine = testConfig.SetupTestRouter()
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

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/client", headers, nil)
	respBodyData := testHttp.ReadResponseBody(resp.Body)

	var parsedRespBody map[string]api_helper.DataResponseObject
	testJson.JSONUnmarshaler(respBodyData, &parsedRespBody)

	v, _ := parsedRespBody["data"].Attributes.(map[string]interface{})

	data := testDB.Scan("oauth2_clients", "data")

	var dataVal map[string]interface{}
	testJson.JSONUnmarshaler(data, &dataVal)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
	assert.Equal(s.T(), v["client_id"], dataVal["ID"])
	assert.Equal(s.T(), v["client_secret"], dataVal["Secret"])
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithInvalidBasicAuth() {
	headers := http.Header{}
	// Basic auth with username = admin and password = password
	headers.Set("Authorization", "Basic invalid")

	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/client", headers, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithoutBasicAuth() {
	resp := testHttp.PerformRequest(s.engine, "POST", "/api/v1/client", nil, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}
