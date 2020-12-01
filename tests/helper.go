package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/foolin/goview/supports/ginview"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gutakk/go-google-scraper/config"
)

// Helper function to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("go-google-scraper", store))

	if withTemplates {
		router.HTMLRender = ginview.New(config.AppGoviewConfig())
		router.Static("/dist", "./dist")
	}

	return router
}

func PerformRequest(r http.Handler, method, path string, headers http.Header, payload url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	req.Header = headers

	response := httptest.NewRecorder()

	r.ServeHTTP(response, req)

	return response
}

func ConstructTestDsn() string {
	host := "localhost"
	port := "5432"
	dbName := "go_google_scraper_test"
	username := "postgres"

	return fmt.Sprintf("sslmode=disable host=%s port=%s dbname=%s user=%s",
		host,
		port,
		dbName,
		username,
	)
}
