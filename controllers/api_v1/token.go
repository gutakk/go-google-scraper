package api_v1

import (
	"net/http"

	"github.com/gutakk/go-google-scraper/oauth"
	"github.com/gutakk/go-google-scraper/services/login_api_service"

	"github.com/gin-gonic/gin"
)

type TokenAPIController struct{}

// TODO: Unit test in login API PR as grant type need to change to password
func (t *TokenAPIController) ApplyRoutes(engine *gin.RouterGroup) {
	engine.POST("/token", t.generateToken)
}

func (t *TokenAPIController) generateToken(c *gin.Context) {
	server := oauth.GetOAuthServer()
	server.SetPasswordAuthorizationHandler(login_api_service.PasswordAuthorizationHandler)

	err := oauth.GetOAuthServer().HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, nil)
	}
}
