package google_scraping_service

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

type GoogleResultParser struct {
	Resp *http.Response
}

type GoogleResult struct {
	Links                  int
	NonAdwords             int
	NonAdwordLinks         []string
	TopPostionAdwords      int
	TopPositionAdwordLinks []string
	TotalAdwords           int
}

func (g *GoogleResultParser) ParseGoogleResponse() (GoogleResult, error) {
	doc, err := goquery.NewDocumentFromReader(g.Resp.Body)
	if err != nil {
		return GoogleResult{}, err
	}

	googleResult := GoogleResult{
		Links:                  g.countLinks(doc),
		NonAdwords:             g.countNonAdwords(doc),
		NonAdwordLinks:         g.fetchNonAdwordLinks(doc),
		TopPostionAdwords:      g.countTopPosAdwords(doc),
		TopPositionAdwordLinks: g.fetchTopPosAdwordLinks(doc),
		TotalAdwords:           g.countTotalAdwords(doc),
	}

	return googleResult, nil
}

func (g *GoogleResultParser) countLinks(doc *goquery.Document) int {
	return len(g.parseLinks(doc, "a"))
}

func (g *GoogleResultParser) countNonAdwords(doc *goquery.Document) int {
	return doc.Find("#rso > div[class=g]").Length()
}

func (g *GoogleResultParser) countTopPosAdwords(doc *goquery.Document) int {
	return doc.Find("#tads > div").Length()
}

func (g *GoogleResultParser) countTotalAdwords(doc *goquery.Document) int {
	return doc.Find("#tadsb > div").Length() + g.countTopPosAdwords(doc)
}

func (g *GoogleResultParser) fetchNonAdwordLinks(doc *goquery.Document) []string {
	return g.parseLinks(doc, "#rso > div[class=g] .yuRUbf > a")
}

func (g *GoogleResultParser) fetchTopPosAdwordLinks(doc *goquery.Document) []string {
	return g.parseLinks(doc, "#tads > div .Krnil")
}

func (g *GoogleResultParser) parseLinks(doc *goquery.Document, selector string) []string {
	var links []string

	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		if exists {
			links = append(links, href)
		}
	})

	return links
}
