package oauth_service_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/models"
	"github.com/gutakk/go-google-scraper/oauth"
	"github.com/gutakk/go-google-scraper/services/oauth_service"
	testDB "github.com/gutakk/go-google-scraper/tests/db"
	"github.com/gutakk/go-google-scraper/tests/oauth_test"
	"github.com/gutakk/go-google-scraper/tests/path_test"

	"github.com/bxcodec/faker/v3"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	gin.SetMode(gin.TestMode)

	err := os.Chdir(path_test.GetRoot())
	if err != nil {
		log.Fatal(err)
	}

	config.LoadEnv()
	err = oauth.SetupOAuthServer()
	if err != nil {
		log.Fatal(err)
	}

	database, err := gorm.Open(postgres.Open(testDB.ConstructTestDsn()), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.GetDB = func() *gorm.DB {
		return database
	}

	err = db.GetDB().AutoMigrate(&models.User{})
	if err != nil {
		log.Fatal(err)
	}
}

type LoginAPIServiceDbTestSuite struct {
	suite.Suite
	user        models.User
	oauthClient oauth_test.OAuthClient
}

func (s *LoginAPIServiceDbTestSuite) SetupTest() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
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
