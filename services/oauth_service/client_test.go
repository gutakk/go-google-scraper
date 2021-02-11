package oauth_service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/oauth"
	"github.com/gutakk/go-google-scraper/services/oauth_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/oauth_test"
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
}

func (s *OAuthControllerDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM oauth2_clients")
}

func TestOAuthControllerDbTestSuite(t *testing.T) {
	suite.Run(t, new(OAuthControllerDbTestSuite))
}

func (s *OAuthControllerDbTestSuite) TestGenerateClient() {
	oauthClient, err := oauth_service.GenerateClient()
	var result oauth_test.OAuthClient
	db.GetDB().Table("oauth2_clients").Select("id", "secret", "domain").Scan(&result)

	assert.Equal(s.T(), oauthClient.ClientID, result.ID)
	assert.Equal(s.T(), oauthClient.ClientSecret, result.Secret)
	assert.Equal(s.T(), fmt.Sprintf("http://localhost:%s", os.Getenv("PORT")), result.Domain)
	assert.Equal(s.T(), nil, err)
}
