package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"

	log "github.com/sirupsen/logrus"
)

func PerformRequest(r http.Handler, method, path string, headers http.Header, payload url.Values) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	if err != nil {
		log.Error(errorHelper.RequestInitializeFailure, err)
	}
	return perform(req, r, headers)
}

func PerformFileUploadRequest(r http.Handler, method, path string, headers http.Header, payload *bytes.Buffer) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		log.Error(errorHelper.RequestInitializeFailure, err)
	}
	return perform(req, r, headers)
}

func perform(req *http.Request, r http.Handler, headers http.Header) *httptest.ResponseRecorder {
	req.Header = headers
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	return response
}
