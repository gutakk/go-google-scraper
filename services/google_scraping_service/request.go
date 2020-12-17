package google_scraping_service

import (
	"net/http"
	"net/url"
)

var Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
	url := "https://www.google.com/search?q=" + url.QueryEscape(keyword)
	client := &http.Client{Transport: transport}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
