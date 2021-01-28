package controllers

import (
	"net/http"
	"os"

	"github.com/gutakk/go-google-scraper/config"
	"github.com/gutakk/go-google-scraper/controllers/api_v1"
	html "github.com/gutakk/go-google-scraper/helpers/html"
	"github.com/gutakk/go-google-scraper/middlewares"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-gonic/gin"
)

const (
	NotFoundTitle = "Not Found"
	NotFoundView  = "not_found"
)

func CombineRoutes(engine *gin.Engine) {
	// Not found
	engine.NoRoute(func(c *gin.Context) {
		html.RenderErrorPage(c, http.StatusNotFound, NotFoundView, NotFoundTitle)
	}, ginview.NewMiddleware(config.ErrorGoviewConfig()))

	// No group
	new(HomeController).applyRoutes(engine)

	// Ensure authenticated user group
	new(KeywordController).applyRoutes(EnsureAuthenticatedUserGroup(engine))
	new(LogoutController).applyRoutes(EnsureAuthenticatedUserGroup(engine))

	// Ensure guest user group
	new(RegisterController).applyRoutes(EnsureGuestUserGroup(engine))
	new(LoginController).applyRoutes(EnsureGuestUserGroup(engine))

	// V1 API group
	v1 := engine.Group("/api/v1")
	// Basic Auth API group
	new(api_v1.OAuthController).ApplyRoutes(BasicAuthAPIGroup(v1))
	// Public API group
	new(api_v1.LoginAPIController).ApplyRoutes(PublicAPIGroup(v1))
	// Private API group
	new(api_v1.KeywordAPIController).ApplyRoutes(PrivateAPIGroup(v1))
}

func BasicAuthAPIGroup(apiVersion *gin.RouterGroup) *gin.RouterGroup {
	return apiVersion.Group("", gin.BasicAuth(gin.Accounts{
		os.Getenv("BASIC_AUTHENTICATION_USERNAME"): os.Getenv("BASIC_AUTHENTICATION_PASSWORD"),
	}))
}

func PublicAPIGroup(apiVersion *gin.RouterGroup) *gin.RouterGroup {
	return apiVersion.Group("")
}

func PrivateAPIGroup(apiVersion *gin.RouterGroup) *gin.RouterGroup {
	privateAPIGroup := apiVersion.Group("")
	privateAPIGroup.Use(middlewares.ValidateToken)

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
