package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
)

func PerformRequest(r http.Handler, method, path string, headers http.Header, payload url.Values) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	return perform(req, r, headers)
}

func PerformFileUploadRequest(r http.Handler, method, path string, headers http.Header, payload *bytes.Buffer) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, payload)
	return perform(req, r, headers)
}

func perform(req *http.Request, r http.Handler, headers http.Header) *httptest.ResponseRecorder {
	req.Header = headers
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	return response
}
