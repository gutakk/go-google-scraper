package google_search_service

import (
	"net/http"
	"testing"

	errorHelper "github.com/gutakk/go-google-scraper/helpers/error_handler"

	"github.com/dnaeon/go-vcr/recorder"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/assert.v1"
)

func TestParserWithValidGoogleResponse(t *testing.T) {
	r, err := recorder.New("tests/fixture/vcr/valid_keyword")
	if err != nil {
		log.Error(errorHelper.RecordInitializeFailure, err)
	}

	url := "https://www.google.com/search?q=AWS"
	client := &http.Client{Transport: r}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(errorHelper.RequestInitializeFailure, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(errorHelper.RequestFailure, err)
	}

	err = r.Stop()
	if err != nil {
		log.Error(errorHelper.RecordStopFailure, err)
	}

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
	r, err := recorder.New("tests/fixture/vcr/invalid_site")
	if err != nil {
		log.Error(errorHelper.RecordInitializeFailure, err)
	}

	url := "https://www.golang.org"
	client := &http.Client{Transport: r}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Error(errorHelper.RequestInitializeFailure, err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(errorHelper.RequestFailure, err)
	}

	err = r.Stop()
	if err != nil {
		log.Error(errorHelper.RecordStopFailure, err)
	}

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
