package google_search_service

import (
	"net/http"
	"net/url"

	"github.com/golang/glog"
)

var Request = func(keyword string, transport http.RoundTripper) (*http.Response, error) {
	url := "https://www.google.com/search?q=" + url.QueryEscape(keyword)
	client := &http.Client{Transport: transport}

	req, requesterErr := http.NewRequest("GET", url, nil)
	if requesterErr != nil {
		glog.Errorf("Cannot init requester: %s", requesterErr)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.88 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
