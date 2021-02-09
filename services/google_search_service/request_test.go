package google_search_service

import (
	"io/ioutil"
	"strings"
	"testing"

	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"

	"github.com/dnaeon/go-vcr/recorder"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/valid_keyword")
	if err != nil {
		log.Error(errorHelper.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("AWS", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorHelper.RecordStopFailure, err)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>AWS") > 0

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/blank_space_keyword")
	if err != nil {
		log.Error(errorHelper.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("A W S", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorHelper.RecordStopFailure, err)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>A W S") > 0

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithThaiKeyword(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/thai_keyword")
	if err != nil {
		log.Error(errorHelper.RecordInitializeFailure, err)
	}

	resp, requestErr := Request("สวัสดี", r)

	err = r.Stop()
	if err != nil {
		log.Error(errorHelper.RecordStopFailure, err)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>สวัสดี") > 0

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}
