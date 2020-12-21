package google_search_service

import (
	"net/http"
	"reflect"

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

type parsingChannel struct {
	countBottomPositionAdwordsCh  chan int
	countLinksCh                  chan int
	countNonAdwordsCh             chan int
	countTopPositionAdwordsCh     chan int
	fetchNonAdwordLinksCh         chan []string
	fetchTopPositionAdwordLinksCh chan []string
}

var ParseGoogleResponse = func(googleResp *http.Response) (ParsingResult, error) {
	parsingCh := parsingChannel{
		countBottomPositionAdwordsCh:  make(chan int),
		countLinksCh:                  make(chan int),
		countNonAdwordsCh:             make(chan int),
		countTopPositionAdwordsCh:     make(chan int),
		fetchNonAdwordLinksCh:         make(chan []string),
		fetchTopPositionAdwordLinksCh: make(chan []string),
	}

	doc, err := goquery.NewDocumentFromReader(googleResp.Body)
	if err != nil {
		return ParsingResult{}, err
	}

	go countBottomPositionAdwords(doc, parsingCh.countBottomPositionAdwordsCh)
	go countLinks(doc, parsingCh.countLinksCh)
	go countNonAdwords(doc, parsingCh.countNonAdwordsCh)
	go countTopPositionAdwords(doc, parsingCh.countTopPositionAdwordsCh)

	go fetchNonAdwordLinks(doc, parsingCh.fetchNonAdwordLinksCh)
	go fetchTopPositionAdwordLinks(doc, parsingCh.fetchTopPositionAdwordLinksCh)

	htmlCode, _ := doc.Html()

	parsingResult := parsingCh.getParsingResultFromChannel()
	parsingResult.HtmlCode = htmlCode

	return parsingResult, nil
}

func (pc *parsingChannel) getParsingResultFromChannel() ParsingResult {
	parsingChannelLength := reflect.TypeOf(*pc).NumField()
	var bottomPositionAdwordsCount, topPositionAdwordsCount, linksCount, nonAdwordsCount int
	var nonAdwordLinks, topPositionAdwordLinks []string

	for i := 0; i < parsingChannelLength; i++ {
		select {
		case val := <-pc.countBottomPositionAdwordsCh:
			bottomPositionAdwordsCount = val
		case val := <-pc.countTopPositionAdwordsCh:
			topPositionAdwordsCount = val
		case val := <-pc.countLinksCh:
			linksCount = val
		case val := <-pc.countNonAdwordsCh:
			nonAdwordsCount = val
		case val := <-pc.fetchNonAdwordLinksCh:
			nonAdwordLinks = val
		case val := <-pc.fetchTopPositionAdwordLinksCh:
			topPositionAdwordLinks = val
		}
	}

	return ParsingResult{
		LinksCount:             linksCount,
		NonAdwordsCount:        nonAdwordsCount,
		NonAdwordLinks:         nonAdwordLinks,
		TopPostionAdwordsCount: topPositionAdwordsCount,
		TopPositionAdwordLinks: topPositionAdwordLinks,
		TotalAdwordsCount:      bottomPositionAdwordsCount + topPositionAdwordsCount,
	}

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
