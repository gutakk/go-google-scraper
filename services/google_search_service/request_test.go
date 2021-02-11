package google_search_service

import (
	"testing"

	errorconf "github.com/gutakk/go-google-scraper/config/error"
	"github.com/gutakk/go-google-scraper/helpers/log"
	testhttp "github.com/gutakk/go-google-scraper/tests/http"

	"github.com/dnaeon/go-vcr/recorder"
	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/valid_keyword")
	if err != nil {
		log.Error(errorconf.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("AWS", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testhttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testhttp.ValidateResponseBody(bodyByte, "<title>AWS")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/blank_space_keyword")
	if err != nil {
		log.Error(errorconf.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("A W S", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testhttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testhttp.ValidateResponseBody(bodyByte, "<title>A W S")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithThaiKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/thai_keyword")
	if err != nil {
		log.Error(errorconf.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("สวัสดี", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorconf.RecordStopFailure, err)
	}

	bodyByte := testhttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testhttp.ValidateResponseBody(bodyByte, "<title>สวัสดี")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}
