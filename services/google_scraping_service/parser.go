package google_scraping_service

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type GoogleResponseParser struct {
	GoogleResponse *http.Response
}

type ScrapingResult struct {
	Links                  int
	NonAdwords             int
	NonAdwordLinks         []string
	TopPostionAdwords      int
	TopPositionAdwordLinks []string
	TotalAdwords           int
}

func (g *GoogleResponseParser) ParseGoogleResponse() (ScrapingResult, error) {
	doc, err := goquery.NewDocumentFromReader(g.GoogleResponse.Body)
	if err != nil {
		return ScrapingResult{}, err
	}

	scrapingResult := ScrapingResult{
		Links:                  g.countLinks(doc),
		NonAdwords:             g.countNonAdwords(doc),
		NonAdwordLinks:         g.fetchNonAdwordLinks(doc),
		TopPostionAdwords:      g.countTopPositionAdwords(doc),
		TopPositionAdwordLinks: g.fetchTopPositionAdwordLinks(doc),
		TotalAdwords:           g.countTotalAdwords(doc),
	}

	return scrapingResult, nil
}

func (g *GoogleResponseParser) countLinks(doc *goquery.Document) int {
	return len(g.parseLinks(doc, "a"))
}

func (g *GoogleResponseParser) countNonAdwords(doc *goquery.Document) int {
	return doc.Find("#rso > div[class=g]").Length()
}

func (g *GoogleResponseParser) countTopPositionAdwords(doc *goquery.Document) int {
	return doc.Find("#tads > div").Length()
}

func (g *GoogleResponseParser) countTotalAdwords(doc *goquery.Document) int {
	return doc.Find("#tadsb > div").Length() + g.countTopPositionAdwords(doc)
}

func (g *GoogleResponseParser) fetchNonAdwordLinks(doc *goquery.Document) []string {
	return g.parseLinks(doc, "#rso > div[class=g] .yuRUbf > a")
}

func (g *GoogleResponseParser) fetchTopPositionAdwordLinks(doc *goquery.Document) []string {
	return g.parseLinks(doc, "#tads > div .Krnil")
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
