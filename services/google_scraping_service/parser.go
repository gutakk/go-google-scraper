package google_scraping_service

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type GoogleResponseParser struct {
	GoogleResponse *http.Response
}

type ParsingResult struct {
	HtmlCode               string
	LinksCount             int
	NonAdwordsCount        int
	NonAdwordLinks         []string
	TopPostionAdwordsCount int
	TopPositionAdwordLinks []string
	TotalAdwordsCount      int
}

var ParseGoogleResponse = func(googleResp *http.Response) (ParsingResult, error) {
	countBottomPositionAdwordsCh := make(chan int)
	countLinksCh := make(chan int)
	countNonAdwordsCh := make(chan int)
	countTopPositionAdwordsCh := make(chan int)

	fetchNonAdwordLinksCh := make(chan []string)
	fetchTopPositionAdwordLinksCh := make(chan []string)

	doc, err := goquery.NewDocumentFromReader(googleResp.Body)
	if err != nil {
		return ParsingResult{}, err
	}

	go countBottomPositionAdwords(doc, countBottomPositionAdwordsCh)
	go countLinks(doc, countLinksCh)
	go countNonAdwords(doc, countNonAdwordsCh)
	go countTopPositionAdwords(doc, countTopPositionAdwordsCh)

	go fetchNonAdwordLinks(doc, fetchNonAdwordLinksCh)
	go fetchTopPositionAdwordLinks(doc, fetchTopPositionAdwordLinksCh)

	bottomPositionAdwordsCount := <-countBottomPositionAdwordsCh
	htmlCode, _ := doc.Html()
	topPositionAdwordsCount := <-countTopPositionAdwordsCh

	parsingResult := ParsingResult{
		HtmlCode:               htmlCode,
		LinksCount:             <-countLinksCh,
		NonAdwordsCount:        <-countNonAdwordsCh,
		NonAdwordLinks:         <-fetchNonAdwordLinksCh,
		TopPostionAdwordsCount: topPositionAdwordsCount,
		TopPositionAdwordLinks: <-fetchTopPositionAdwordLinksCh,
		TotalAdwordsCount:      bottomPositionAdwordsCount + topPositionAdwordsCount,
	}

	return parsingResult, nil
}

func countBottomPositionAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#tadsb > div").Length()
}

func countLinks(doc *goquery.Document, ch chan int) {
	ch <- len(parseLinks(doc, "a"))
}

func countNonAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#rso .yuRUbf").Length()
}

func countTopPositionAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#tads > div").Length()
}

func fetchNonAdwordLinks(doc *goquery.Document, ch chan []string) {
	ch <- parseLinks(doc, "#rso .yuRUbf > a")
}

func fetchTopPositionAdwordLinks(doc *goquery.Document, ch chan []string) {
	ch <- parseLinks(doc, "#tads > div .Krnil")
}

func parseLinks(doc *goquery.Document, selector string) []string {
	var links []string

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	return links
}
