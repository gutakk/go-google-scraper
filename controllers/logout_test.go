package controllers_test

import (
	"net/http"
	"testing"

	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	"github.com/gutakk/go-google-scraper/tests/fabricator"
	"github.com/gutakk/go-google-scraper/tests/fixture"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/suite"
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
	user := fabricator.FabricateUser("test@email.com", "testPassword")

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
