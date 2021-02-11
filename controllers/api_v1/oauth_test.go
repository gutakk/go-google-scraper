package api_v1_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/api_helper"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

func init() {
	gin.SetMode(gin.TestMode)

	path_test.ChangeToRootDir()

	config.LoadEnv()

	err := oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(errorconf.StartOAuthServerFailure, err)
	}

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
	respBodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	var parsedRespBody map[string]api_helper.DataResponseObject
	err = json.Unmarshal(respBodyData, &parsedRespBody)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

	v, _ := parsedRespBody["data"].Attributes.(map[string]interface{})

	var data []byte
	row := db.GetDB().Table("oauth2_clients").Select("data").Row()
	err = row.Scan(&data)
	if err != nil {
		log.Error(errorconf.ScanRowFailure, err)
	}

	var dataVal map[string]interface{}
	err = json.Unmarshal(data, &dataVal)
	if err != nil {
		log.Error(errorconf.JSONUnmarshalFailure, err)
	}

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
