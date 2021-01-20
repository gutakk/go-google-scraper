package google_search_service

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/golang/glog"
	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	r, recorderErr := recorder.New("tests/fixture/vcr/valid_keyword")
	if recorderErr != nil {
		glog.Errorf("Cannot init recorder: %s", recorderErr)
	}

	resp, requestErr := Request("AWS", r)
	if requestErr != nil {
		glog.Errorf("Cannot make a request: %s", requestErr)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>AWS") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	stopRecorderErr := r.Stop()
	if stopRecorderErr != nil {
		glog.Errorf("Cannot stop the recorder: %s", stopRecorderErr)
	}
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	r, recorderErr := recorder.New("tests/fixture/vcr/blank_space_keyword")
	if recorderErr != nil {
		glog.Errorf("Cannot init recorder: %s", recorderErr)
	}

	resp, requestErr := Request("A W S", r)
	if requestErr != nil {
		glog.Errorf("Cannot make a request: %s", requestErr)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>A W S") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	stopRecorderErr := r.Stop()
	if stopRecorderErr != nil {
		glog.Errorf("Cannot stop the recorder: %s", stopRecorderErr)
	}
}

func TestRequestWithThaiKeyword(t *testing.T) {
	r, recorderErr := recorder.New("tests/fixture/vcr/thai_keyword")
	if recorderErr != nil {
		glog.Errorf("Cannot init recorder: %s", recorderErr)
	}

	resp, requestErr := Request("สวัสดี", r)
	if requestErr != nil {
		glog.Errorf("Cannot make a request: %s", requestErr)
	}

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>สวัสดี") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	stopRecorderErr := r.Stop()
	if stopRecorderErr != nil {
		glog.Errorf("Cannot stop the recorder: %s", stopRecorderErr)
	}
}
