package tests

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"

	"github.com/dnaeon/go-vcr/recorder"
)

func PerformRequest(r http.Handler, method, path string, headers http.Header, payload url.Values) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, strings.NewReader(payload.Encode()))
	if err != nil {
		log.Error(errorconf.RequestInitializeFailure, err)
	}
	return perform(req, r, headers)
}

func PerformFileUploadRequest(r http.Handler, method, path string, headers http.Header, payload *bytes.Buffer) *httptest.ResponseRecorder {
	req, err := http.NewRequest(method, path, payload)
	if err != nil {
		log.Error(errorconf.RequestInitializeFailure, err)
	}
	return perform(req, r, headers)
}

func perform(req *http.Request, r http.Handler, headers http.Header) *httptest.ResponseRecorder {
	req.Header = headers
	response := httptest.NewRecorder()
	r.ServeHTTP(response, req)
	return response
}

func ReadResponseBody(respBody interface{}) []byte {
	var bodyByte []byte
	var err error
	byteBufferBody, isByteBuffer := respBody.(*bytes.Buffer)
	readCloserBody, isReadCloser := respBody.(io.ReadCloser)

	if isByteBuffer {
		bodyByte, err = ioutil.ReadAll(byteBufferBody)
	} else if isReadCloser {
		bodyByte, err = ioutil.ReadAll(readCloserBody)
	}

	if err != nil {
		log.Error(errorconf.ReadResponseBodyFailure, err)
	}

	return bodyByte
}

func ValidateResponseBody(bodyByte []byte, expected string) bool {
	return strings.Index(string(bodyByte), expected) > 0
}

func NewRecorder(cassetteName string) *recorder.Recorder {
	recorder, err := recorder.New(cassetteName)
	if err != nil {
		log.Error(errorconf.RecorderInitializeFailure, err)
	}

	return recorder
}

func StopRecorder(recorder *recorder.Recorder) {
	err := recorder.Stop()
	if err != nil {
		log.Error(errorconf.RecorderStopFailure, err)
	}
}
