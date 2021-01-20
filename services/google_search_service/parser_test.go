package google_search_service

import (
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"github.com/golang/glog"
	"gopkg.in/go-playground/assert.v1"
)

func TestParserWithValidGoogleResponse(t *testing.T) {
	r, recorderErr := recorder.New("tests/fixture/vcr/valid_keyword")
	if recorderErr != nil {
		glog.Errorf("Cannot init recorder: %s", recorderErr)
	}

	url := "https://www.google.com/search?q=AWS"
	client := &http.Client{Transport: r}
	req, requesterErr := http.NewRequest("GET", url, nil)
	if requesterErr != nil {
		glog.Errorf("Cannot init requester: %s", requesterErr)
	}

	resp, requestErr := client.Do(req)
	if requestErr != nil {
		glog.Errorf("Cannot make a request: %s", requestErr)
	}

	stopRecorderErr := r.Stop()
	if stopRecorderErr != nil {
		glog.Errorf("Cannot stop the recorder: %s", stopRecorderErr)
	}

	parsingResult, err := ParseGoogleResponse(resp)

	assert.Equal(t, nil, err)
	assert.Equal(t, 90, parsingResult.LinksCount)
	assert.Equal(t, 7, parsingResult.NonAdwordsCount)
	assert.Equal(t, 4, parsingResult.TopPostionAdwordsCount)
	assert.Equal(t, 4, parsingResult.TotalAdwordsCount)
	assert.NotEqual(t, nil, parsingResult.HtmlCode)
	assert.NotEqual(t, nil, parsingResult.NonAdwordLinks)
	assert.NotEqual(t, nil, parsingResult.TopPositionAdwordLinks)
}

func TestParserWithNotGoogleSearchPage(t *testing.T) {
	r, recorderErr := recorder.New("tests/fixture/vcr/invalid_site")
	if recorderErr != nil {
		glog.Errorf("Cannot init recorder: %s", recorderErr)
	}

	url := "https://www.golang.org"
	client := &http.Client{Transport: r}
	req, requesterErr := http.NewRequest("GET", url, nil)
	if requesterErr != nil {
		glog.Errorf("Cannot init requester: %s", requesterErr)
	}

	resp, requestErr := client.Do(req)
	if requestErr != nil {
		glog.Errorf("Cannot make a request: %s", requestErr)
	}

	stopRecorderErr := r.Stop()
	if stopRecorderErr != nil {
		glog.Errorf("Cannot stop the recorder: %s", stopRecorderErr)
	}

	parsingResult, err := ParseGoogleResponse(resp)

	assert.Equal(t, nil, err)
	assert.Equal(t, 16, parsingResult.LinksCount)
	assert.Equal(t, 0, parsingResult.NonAdwordsCount)
	assert.Equal(t, 0, parsingResult.TopPostionAdwordsCount)
	assert.Equal(t, 0, parsingResult.TotalAdwordsCount)
	assert.Equal(t, nil, parsingResult.NonAdwordLinks)
	assert.Equal(t, nil, parsingResult.TopPositionAdwordLinks)
	assert.NotEqual(t, nil, parsingResult.HtmlCode)
}
