package google_scraping_service

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	r, _ := recorder.New("../../tests/fixture/vcr/valid_keyword")

	googleRequest := GoogleRequest{Keyword: "AWS", Transport: r}
	resp, _ := googleRequest.Request()

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>AWS") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	_ = r.Stop()
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	r, _ := recorder.New("../../tests/fixture/vcr/blank_space_keyword")

	googleRequest := GoogleRequest{Keyword: "A W S", Transport: r}
	resp, _ := googleRequest.Request()

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>A W S") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	_ = r.Stop()
}

func TestRequestWithThaiKeyword(t *testing.T) {
	r, _ := recorder.New("../../tests/fixture/vcr/thai_keyword")

	googleRequest := GoogleRequest{Keyword: "สวัสดี", Transport: r}
	resp, _ := googleRequest.Request()

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>สวัสดี") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)

	_ = r.Stop()
}
