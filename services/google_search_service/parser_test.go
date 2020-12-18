package google_search_service

import (
	"net/http"
	"testing"

	"github.com/dnaeon/go-vcr/recorder"
	"gopkg.in/go-playground/assert.v1"
)

func TestParserWithValidGoogleResponse(t *testing.T) {
	r, _ := recorder.New("tests/fixture/vcr/valid_keyword")

	url := "https://www.google.com/search?q=AWS"
	client := &http.Client{Transport: r}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	_ = r.Stop()

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
	r, _ := recorder.New("tests/fixture/vcr/invalid_site")

	url := "https://www.golang.org"
	client := &http.Client{Transport: r}
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := client.Do(req)

	_ = r.Stop()

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
