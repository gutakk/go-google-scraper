package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func SetupMiddlewares(router *gin.Engine) *gin.Engine {
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("go-google-scraper", store))
	router.Use(CurrentUser)

	return router
}
