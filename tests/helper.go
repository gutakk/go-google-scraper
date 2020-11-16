package tests

import (
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// Helper function to create a router during testing
func GetRouter(withTemplates bool) *gin.Engine {
	router := gin.Default()

	if withTemplates {
		router.LoadHTMLGlob("templates/*")
		router.Static("/dist", "./dist")
	}

	return router
}

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	return w
}
