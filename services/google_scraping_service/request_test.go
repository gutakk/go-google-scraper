package google_scraping_service

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"gopkg.in/go-playground/assert.v1"
)

func TestRequest(t *testing.T) {
	r, _ := recorder.New("../../tests/fixture/vcr")
	defer r.Stop()

	googleRequest := GoogleRequest{Keyword: "AWS", Transport: r}
	resp, _ := googleRequest.Request()

	p, err := ioutil.ReadAll(resp.Body)
	isGoogleSearchPage := err == nil && strings.Index(string(p), "<title>AWS") > 0

	assert.Equal(t, nil, err)
	assert.Equal(t, true, isGoogleSearchPage)
}
