package controllers_test

import (
	"net/http"
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/db"
	"github.com/gutakk/go-google-scraper/helpers/log"
	"github.com/gutakk/go-google-scraper/models"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/assert.v1"
)

type LogoutTestSuite struct {
	suite.Suite
	engine *gin.Engine
}

func (s *LogoutTestSuite) SetupTest() {
	s.engine = testConfig.SetupTestRouter()
}

func TestLogoutTestSuit(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}

func (s *LogoutTestSuite) TestLogoutWithAuthenticatedUser() {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("testPassword"), bcrypt.DefaultCost)
	if err != nil {
		log.Error(errorconf.HashPasswordFailure, err)
	}
	user := models.User{Email: "test@email.com", Password: string(hashedPassword)}
	db.GetDB().Create(&user)

	cookie := fixture.GenerateCookie("user_id", user.ID)
	headers := http.Header{}
	headers.Set("Cookie", cookie.Name+"="+cookie.Value)

	response := testHttp.PerformRequest(s.engine, "POST", "/logout", headers, nil)

	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/", response.Header().Get("Location"))
}

func (s *LogoutTestSuite) TestLogoutWithGuestUser() {
	response := testHttp.PerformRequest(s.engine, "POST", "/logout", nil, nil)

	// TODO: Research the flash messge assertion solution
	assert.Equal(s.T(), http.StatusFound, response.Code)
	assert.Equal(s.T(), "/login", response.Header().Get("Location"))
}
