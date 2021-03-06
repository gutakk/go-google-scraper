package google_search_service

import (
	"testing"

	testhttp "github.com/gutakk/go-google-scraper/tests/http"

	"gopkg.in/go-playground/assert.v1"
)

func TestParserWithValidGoogleResponse(t *testing.T) {
	recorder := testhttp.NewRecorder("tests/fixture/vcr/valid_keyword")

	url := "https://www.google.com/search?q=AWS"
	resp := testhttp.PerformClientRequest(url, recorder)

	testhttp.StopRecorder(recorder)

	parsingResult, parsingError := ParseGoogleResponse(resp)

	assert.Equal(t, nil, parsingError)
	assert.Equal(t, 90, parsingResult.LinksCount)
	assert.Equal(t, 7, parsingResult.NonAdwordsCount)
	assert.Equal(t, 4, parsingResult.TopPostionAdwordsCount)
	assert.Equal(t, 4, parsingResult.TotalAdwordsCount)
	assert.NotEqual(t, nil, parsingResult.HtmlCode)
	assert.NotEqual(t, nil, parsingResult.NonAdwordLinks)
	assert.NotEqual(t, nil, parsingResult.TopPositionAdwordLinks)
}

func TestParserWithNotGoogleSearchPage(t *testing.T) {
	recorder := testhttp.NewRecorder("tests/fixture/vcr/invalid_site")

	url := "https://www.golang.org"
	resp := testhttp.PerformClientRequest(url, recorder)

	testhttp.StopRecorder(recorder)

	parsingResult, parsingErr := ParseGoogleResponse(resp)

	assert.Equal(t, nil, parsingErr)
	assert.Equal(t, 16, parsingResult.LinksCount)
	assert.Equal(t, 0, parsingResult.NonAdwordsCount)
	assert.Equal(t, 0, parsingResult.TopPostionAdwordsCount)
	assert.Equal(t, 0, parsingResult.TotalAdwordsCount)
	assert.Equal(t, nil, parsingResult.NonAdwordLinks)
	assert.Equal(t, nil, parsingResult.TopPositionAdwordLinks)
	assert.NotEqual(t, nil, parsingResult.HtmlCode)
}
