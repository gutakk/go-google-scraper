package google_scraping_service

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type GoogleResponseParser struct {
	GoogleResponse *http.Response
}

type ScrapingResult struct {
	LinksCount             int
	NonAdwordsCount        int
	NonAdwordLinks         []string
	TopPostionAdwordsCount int
	TopPositionAdwordLinks []string
	TotalAdwordsCount      int
}

func (g *GoogleResponseParser) ParseGoogleResponse() (ScrapingResult, error) {
	countBottomPositionAdwordsCh := make(chan int)
	countLinksCh := make(chan int)
	countNonAdwordsCh := make(chan int)
	countTopPositionAdwordsCh := make(chan int)

	fetchNonAdwordLinksCh := make(chan []string)
	fetchTopPositionAdwordLinksCh := make(chan []string)

	doc, err := goquery.NewDocumentFromReader(g.GoogleResponse.Body)
	if err != nil {
		return ScrapingResult{}, err
	}

	go g.countBottomPositionAdwords(doc, countBottomPositionAdwordsCh)
	go g.countLinks(doc, countLinksCh)
	go g.countNonAdwords(doc, countNonAdwordsCh)
	go g.countTopPositionAdwords(doc, countTopPositionAdwordsCh)

	go g.fetchNonAdwordLinks(doc, fetchNonAdwordLinksCh)
	go g.fetchTopPositionAdwordLinks(doc, fetchTopPositionAdwordLinksCh)

	bottomPositionAdwordsCount := <-countBottomPositionAdwordsCh
	topPositionAdwordsCount := <-countTopPositionAdwordsCh

	scrapingResult := ScrapingResult{
		LinksCount:             <-countLinksCh,
		NonAdwordsCount:        <-countNonAdwordsCh,
		NonAdwordLinks:         <-fetchNonAdwordLinksCh,
		TopPostionAdwordsCount: topPositionAdwordsCount,
		TopPositionAdwordLinks: <-fetchTopPositionAdwordLinksCh,
		TotalAdwordsCount:      bottomPositionAdwordsCount + topPositionAdwordsCount,
	}
	return scrapingResult, nil
}

func (g *GoogleResponseParser) countLinks(doc *goquery.Document, ch chan int) {
	ch <- len(g.parseLinks(doc, "a"))
}

func (g *GoogleResponseParser) countNonAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#rso > div[class=g]").Length()
}

func (g *GoogleResponseParser) countTopPositionAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#tads > div").Length()
}

func (g *GoogleResponseParser) countBottomPositionAdwords(doc *goquery.Document, ch chan int) {
	ch <- doc.Find("#tadsb > div").Length()
}

func (g *GoogleResponseParser) fetchNonAdwordLinks(doc *goquery.Document, ch chan []string) {
	ch <- g.parseLinks(doc, "#rso > div[class=g] .yuRUbf > a")
}

func (g *GoogleResponseParser) fetchTopPositionAdwordLinks(doc *goquery.Document, ch chan []string) {
	ch <- g.parseLinks(doc, "#tads > div .Krnil")
}

func (g *GoogleResponseParser) parseLinks(doc *goquery.Document, selector string) []string {
	var links []string

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	return links
}
