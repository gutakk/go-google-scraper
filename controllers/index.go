package controllers

import (
	"net/http"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers/api"
	"github.com/gutakk/go-google-scraper/middlewares"
)

func CombineRoutes(engine *gin.Engine) {
	// No group
	new(HomeController).applyRoutes(engine)

	// Ensure authenticated user group
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))
	new(LogoutController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Ensure guest user group
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))

	// Basic Auth API group
	new(api.OAuthController).ApplyRoutes(BasicAuthAPIGroup(engine))

	// Public API group
	new(api.TokenAPIController).ApplyRoutes(PublicAPIGroup(engine))

	new(api.DummyAPIController).ApplyRoutes(PrivateAPIGroup(engine))
}

func BasicAuthAPIGroup(engine *gin.Engine) *gin.RouterGroup {
	return engine.Group("/api", gin.BasicAuth(gin.Accounts{
		"foo": "bar",
	}))
}

func PublicAPIGroup(engine *gin.Engine) *gin.RouterGroup {
	return engine.Group("/api")
}

func PrivateAPIGroup(engine *gin.Engine) *gin.RouterGroup {
	privateAPIGroup := engine.Group("/api")
	privateAPIGroup.Use(validateToken)

	return privateAPIGroup
}

func EnsureAuthenticatedUserGroup(engine *gin.Engine) *gin.RouterGroup {
	ensureAuthenticatedUserGroup := engine.Group("")
	ensureAuthenticatedUserGroup.Use(middlewares.EnsureAuthenticatedUser)

	return ensureAuthenticatedUserGroup
}

func EnsureGuestUserGroup(engine *gin.Engine) *gin.RouterGroup {
	mw := ginview.NewMiddleware(config.AuthenticationGoviewConfig())

	ensureGuestUserGroup := engine.Group("", mw)
	ensureGuestUserGroup.Use(middlewares.EnsureGuestUser)

	return ensureGuestUserGroup
}

func validateToken(c *gin.Context) {
	_, err := config.GetOAuthServer().ValidationBearerToken(c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusUnauthorized)
		c.Abort()
	}
}
