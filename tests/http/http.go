package tests

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func PerformRequest(r http.Handler, method, path string, headers http.Header, payload url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	req.Header = headers

	response := httptest.NewRecorder()

	r.ServeHTTP(response, req)

	return response
}
