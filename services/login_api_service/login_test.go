package login_api_service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	"github.com/gutakk/go-google-scraper/services/login_api_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
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

type LoginAPIServiceDbTestSuite struct {
	suite.Suite
	user        models.User
	oauthClient oauth_test.OAuthClient
}

func (s *LoginAPIServiceDbTestSuite) SetupTest() {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
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
	userID, err := login_api_service.PasswordAuthorizationHandler(s.user.Email, "password")

	assert.Equal(s.T(), fmt.Sprint(s.user.ID), userID)
	assert.Equal(s.T(), nil, err)
}

func (s *LoginAPIServiceDbTestSuite) TestPasswordAuthorizationHandlerInvalidEmail() {
	userID, err := login_api_service.PasswordAuthorizationHandler("invalidEmail", "password")

	assert.Equal(s.T(), "", userID)
	assert.Equal(s.T(), errors.ErrInvalidClient, err)
}

func (s *LoginAPIServiceDbTestSuite) TestPasswordAuthorizationHandlerInvalidPassword() {
	userID, err := login_api_service.PasswordAuthorizationHandler(s.user.Email, "invalidPassword")

	assert.Equal(s.T(), "", userID)
	assert.Equal(s.T(), errors.ErrInvalidClient, err)
}
