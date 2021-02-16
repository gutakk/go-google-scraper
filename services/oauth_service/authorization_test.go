package oauth_service_test

import (
	"fmt"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/services/oauth_service"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/fabricator"
	testOauth "github.com/gutakk/go-google-scraper/tests/oauth_test"
	testPath "github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
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

type LoginAPIServiceDbTestSuite struct {
	suite.Suite
	user        models.User
	oauthClient testOauth.OAuthClient
}

func (s *LoginAPIServiceDbTestSuite) SetupTest() {
	user := fabricator.FabricateUser(faker.Email(), "password")

	s.oauthClient = testOauth.OAuthClient{
		ID:     "client-id",
		Secret: "client-secret",
		Domain: "http://localhost:8080",
	}
	db.GetDB().Raw("INSERT INTO oauth2_clients(id, secret, domain) VALUES(?, ?, ?)",
		s.oauthClient.ID,
		s.oauthClient.Secret,
		s.oauthClient.Domain,
	)

	s.user = user
}

func (s *LoginAPIServiceDbTestSuite) TearDownTest() {
	db.GetDB().Exec("DELETE FROM users")
	db.GetDB().Exec("DELETE FROM oauth2_clients")
}

func TestLoginAPIServiceDbTestSuite(t *testing.T) {
	suite.Run(t, new(LoginAPIServiceDbTestSuite))
}

func (s *LoginAPIServiceDbTestSuite) TestPasswordAuthorizationHandlerWithValidParams() {
	userID, err := oauth_service.PasswordAuthorizationHandler(s.user.Email, "password")

	assert.Equal(s.T(), fmt.Sprint(s.user.ID), userID)
	assert.Equal(s.T(), nil, err)
}

func (s *LoginAPIServiceDbTestSuite) TestPasswordAuthorizationHandlerInvalidEmail() {
	userID, err := oauth_service.PasswordAuthorizationHandler("invalidEmail", "password")

	assert.Equal(s.T(), "", userID)
	assert.Equal(s.T(), errors.ErrInvalidClient, err)
}

func (s *LoginAPIServiceDbTestSuite) TestPasswordAuthorizationHandlerInvalidPassword() {
	userID, err := oauth_service.PasswordAuthorizationHandler(s.user.Email, "invalidPassword")

	assert.Equal(s.T(), "", userID)
	assert.Equal(s.T(), errors.ErrInvalidClient, err)
}
