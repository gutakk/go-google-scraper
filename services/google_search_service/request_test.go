package google_search_service

import (
	"testing"

	testHttp "github.com/gutakk/go-google-scraper/tests/http"

	"gopkg.in/go-playground/assert.v1"
)

func TestRequestWithValidKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/valid_keyword")

	resp, requestErr := Request("AWS", recorder)

	testHttp.StopRecorder(recorder)

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>AWS")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithBlankSpaceKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/blank_space_keyword")

	resp, requestErr := Request("A W S", recorder)

	testHttp.StopRecorder(recorder)

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>A W S")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}

func TestRequestWithThaiKeyword(t *testing.T) {
	recorder := testHttp.NewRecorder("tests/fixture/vcr/thai_keyword")

	resp, requestErr := Request("สวัสดี", recorder)

	testHttp.StopRecorder(recorder)

	bodyByte := testHttp.ReadResponseBody(resp.Body)
	isGoogleSearchPage := testHttp.ValidateResponseBody(bodyByte, "<title>สวัสดี")

	assert.Equal(t, nil, requestErr)
	assert.Equal(t, true, isGoogleSearchPage)
}
