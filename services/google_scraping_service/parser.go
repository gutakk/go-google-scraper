package google_scraping_service

import (
	"log"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func GoogleResultParser(resp *http.Response) (string, error) {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("@@@@@@@@@@@@@@@@@@@ %v", err)
	}
	log.Printf("===================== %v", countTopPosAdwords(doc))
	log.Printf("@@@@@@@@@@@@@@@@@@@@@ %v", fetchTopPosAdwordLinks(doc))
	log.Printf("##################### %v", countTotalAdwords(doc, 4))
	log.Printf("!!!!!!!!!!!!!!!!!!!!! %v", countNonAdwords(doc))
	log.Printf("$$$$$$$$$$$$$$$$$$$$$ %v", fetchNonAdwordLinks(doc))
	log.Printf("^^^^^^^^^^^^^^^^^^^^^ %v", countLinks(doc))
	return "hello", nil
}

func countLinks(doc *goquery.Document) int {
	return len(parseLinks(doc, "a"))
}

func countNonAdwords(doc *goquery.Document) int {
	return doc.Find("#rso > div[class=g]").Length()
}

func countTopPosAdwords(doc *goquery.Document) int {
	return doc.Find("#tads > div").Length()
}

func countTotalAdwords(doc *goquery.Document, topPosAdwords int) int {
	return doc.Find("#tadsb > div").Length() + topPosAdwords
}

func fetchNonAdwordLinks(doc *goquery.Document) []string {
	return parseLinks(doc, "#rso > div[class=g] .yuRUbf > a")
}

func fetchTopPosAdwordLinks(doc *goquery.Document) []string {
	return parseLinks(doc, "#tads > div .Krnil")
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
