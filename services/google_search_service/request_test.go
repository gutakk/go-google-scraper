package google_search_service

import (
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/valid_keyword")

	resp, requestErr := Request("AWS", recorder)

	err := recorder.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>AWS")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/blank_space_keyword")

	resp, requestErr := Request("A W S", recorder)

	err := recorder.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>A W S")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithThaiKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/thai_keyword")

	resp, requestErr := Request("สวัสดี", recorder)

	err := recorder.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>สวัสดี")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}
