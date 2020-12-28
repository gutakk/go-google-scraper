package api_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers"
	"github.com/gutakk/go-google-scraper/controllers/api"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/oauth"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/gutakk/go-google-scraper/tests/path_test"
	"github.com/stretchr/testify/suite"

	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type OAuthClientResult struct {
	ID     string
	Secret string
	Domain string
}

func init() {
	gin.SetMode(gin.TestMode)

	if err := os.Chdir(path_test.GetRoot()); err != nil {
		panic(err)
	}

	config.LoadEnv()
	oauth.SetupOAuthServer()
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
	new(api.OAuthController).ApplyRoutes(controllers.BasicAuthAPIGroup(s.engine))
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
	data, _ := ioutil.ReadAll(resp.Body)
	var respBody map[string]string
	_ = json.Unmarshal(data, &respBody)

	assert.Equal(s.T(), http.StatusCreated, resp.Code)
	assert.NotEqual(s.T(), nil, respBody["CLIENT_ID"])
	assert.NotEqual(s.T(), nil, respBody["CLIENT_SECRET"])

	var result OAuthClientResult
	db.GetDB().Table("oauth2_clients").Select("id", "secret", "domain").Scan(&result)

	assert.Equal(s.T(), respBody["CLIENT_ID"], result.ID)
	assert.Equal(s.T(), respBody["CLIENT_SECRET"], result.Secret)
	assert.Equal(s.T(), fmt.Sprintf("http://localhost:%s", os.Getenv("APP_PORT")), result.Domain)
}

func (s *OAuthControllerDbTestSuite) TestGenerateClientWithoutBasicAuth() {
	resp := testHttp.PerformRequest(s.engine, "POST", "/api/client", nil, nil)

	assert.Equal(s.T(), http.StatusUnauthorized, resp.Code)
}
