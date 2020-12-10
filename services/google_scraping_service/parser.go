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

func (g *GoogleResponseParser) ParseGoogleResponse() (ParsingResult, error) {

	doc, err := goquery.NewDocumentFromReader(g.GoogleResponse.Body)
	if err != nil {
		return ParsingResult{}, err
	}

	bottomPoistionAdwordsCount := g.countBottomPositionAdwords(doc)
	topPositionAdwordsCount := g.countTopPositionAdwords(doc)

	parsingResult := ParsingResult{
		HtmlCode:               doc.Text(),
		LinksCount:             g.countLinks(doc),
		NonAdwordsCount:        g.countNonAdwords(doc),
		NonAdwordLinks:         g.fetchNonAdwordLinks(doc),
		TopPostionAdwordsCount: topPositionAdwordsCount,
		TopPositionAdwordLinks: g.fetchTopPositionAdwordLinks(doc),
		TotalAdwordsCount:      bottomPoistionAdwordsCount + topPositionAdwordsCount,
	}

	return parsingResult, nil
}

func (g *GoogleResponseParser) countBottomPositionAdwords(doc *goquery.Document) int {
	return doc.Find("#tadsb > div").Length()
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
