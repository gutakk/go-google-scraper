package tests

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type Header struct {
	Key   string
	Value string
}

// Helper function to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	router := gin.Default()
	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	if withTemplates {
		router.LoadHTMLGlob("templates/*")
		router.Static("/dist", "./dist")
	}

	return router
}

func PerformRequest(r http.Handler, method, path string, headers []Header, payload url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	for _, h := range headers {
		req.Header.Add(h.Key, h.Value)
	}

	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
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
