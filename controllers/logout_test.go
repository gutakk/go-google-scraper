package controllers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	testConfig "github.com/gutakk/go-google-scraper/tests/config"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

type LogoutTestSuite struct {
	suite.Suite
	engine *gin.Engine
}

func (s *LogoutTestSuite) SetupTest() {
	s.engine = testConfig.GetRouter(false)
	new(LogoutController).applyRoutes(EnsureAuthenticatedUserGroup(s.engine))
}

func TestLogoutTestSuit(t *testing.T) {
	suite.Run(t, new(LogoutTestSuite))
}

func (s *LogoutTestSuite) TestLogoutWithAuthenticatedUser() {
	cookie := "go-google-scraper=MTYwNjQ2Mjk3MXxEdi1CQkFFQ180SUFBUkFCRUFBQUlmLUNBQUVHYzNSeWFXNW5EQWtBQjNWelpYSmZhV1FFZFdsdWRBWUVBUDRFdFE9PXzl6APqAQw3gAQqlHoXMYrPpnqPFkEP8SRHJZEpl-_LDQ=="
	headers := http.Header{}
	headers.Set("Cookie", cookie)

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
