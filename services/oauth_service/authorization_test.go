package oauth_service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	"github.com/gutakk/go-google-scraper/services/oauth_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/oauth_test"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
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
		log.Fatal(errorconf.StartOAuthServerFailure, err)
	}

	testDB.SetupTestDatabase()
}

type LoginAPIServiceDbTestSuite struct {
	suite.Suite
	user        models.User
	oauthClient oauth_test.OAuthClient
}

func (s *LoginAPIServiceDbTestSuite) SetupTest() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}

	user := models.User{Email: faker.Email(), Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	s.oauthClient = oauth_test.OAuthClient{
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
