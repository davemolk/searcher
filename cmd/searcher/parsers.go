package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// makeParseData returns a struct of selectors for each of the
// search engines (ask, bing, brave, duck duck go, and yahoo).
func (s *searcher) makeParseData() []parseData {
	var pdSlice []parseData

	ask := parseData{
		blurbSelector: "div.PartialSearchResults-item p",
		itemSelector:  "div.PartialSearchResults-item",
		linkSelector:  "a.PartialSearchResults-item-title-link",
		name:          "ask",
	}

	bing := parseData{
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name:          "bing",
	}

	brave := parseData{
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name:          "brave",
	}

	duck := parseData{
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}

	yahoo := parseData{
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
		name:          "yahoo",
	}

	pdSlice = append(pdSlice, ask, bing, brave, duck, yahoo)

	return pdSlice
}

// getAndParseData creates query strings for each search engine,
// makes the appropriate GET request, and parses the results,
// saving each URL and associated blurb.
func (s *searcher) getAndParseData() {
	pdSlice := s.makeParseData()
	chans := s.makeSearchURLs()

	var wg sync.WaitGroup
	tokens := make(chan struct{}, s.config.concurrency)
	for i, ch := range chans {
		for u := range ch {
			wg.Add(1)
			tokens <- struct{}{}
			go func(u string, i int) {
				defer wg.Done()
				// splits into URL and the search term(s)
				urlTerm := strings.Split(u, "GETTERM")
				if len(urlTerm) == 1 {
					urlTerm = append(urlTerm, "") // this will get stripped later
				}
				buf, err := s.makeRequest(urlTerm[0], s.config.timeout)
				if err != nil {
					if s.config.verbose {
						s.errorLog.Printf("error in makeRequest for %s: %v\n", urlTerm[0], err)
					}
					<-tokens
					return
				}
				<-tokens
				s.parseSearchResults(buf, urlTerm[1], pdSlice[i])
			}(u, i)
		}
	}

	wg.Wait()
}

// parseSearchResults parses the response body for each search result,
// saving the URL and associated blurb.
func (s *searcher) parseSearchResults(data *bytes.Buffer, term string, pd parseData) {
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		if s.config.verbose {
			s.errorLog.Printf("goquery error for %s: %v\n", pd.name, err)
		}
		return
	}

	doc.Find(pd.itemSelector).Each(func(i int, g *goquery.Selection) {
		// no link, no point in getting blurb
		link, ok := g.Find(pd.linkSelector).Attr("href")
		if !ok {
			return
		}
		blurb := g.Find(pd.blurbSelector).Text()
		cleanedLink := s.cleanLinks(link)
		cleanedBlurb := s.cleanBlurb(blurb)
		if !s.config.json {
			s.printStdout(cleanedBlurb, cleanedLink)
		}
		if term != "" {
			s.searches.storeSearches(term, cleanedLink, cleanedBlurb)
		} else {
			s.searches.storeSearch(cleanedLink, cleanedBlurb)
		}
	})
}

// cleanBlurb does a bit of tidying up of each input blurb string.
func (s *searcher) cleanBlurb(str string) string {
	cleanB := s.noBlank.ReplaceAllString(str, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

// cleanLinks does a bit of tidying up of each input URL string.
// bing will sometimes encode the links and I haven't bothered to work
// out how to decode them. Maybe one day...
func (s *searcher) cleanLinks(str string) string {
	u, err := url.QueryUnescape(str)
	if err != nil {
		if s.config.verbose {
			s.errorLog.Printf("unable to clean %s: %v\n", str, err)
		}
		return str
	}
	if strings.HasPrefix(u, "//duck") {
		// ddg links will sometimes take the following format:
		// //duckduckgo.com/l/?uddg=actualURLHere/&rut=otherStuff
		removePrefix := strings.Split(u, "=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "&rut")
		u = removeSuffix[0]
	}
	if strings.HasPrefix(u, "https://r.search.yahoo.com/") {
		removePrefix := strings.Split(u, "/RU=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "/RK=")
		u = removeSuffix[0]
	}
	return u
}

// printStdout truncates any blurb with a length longer
// than 200 and prints to stdout.
func (s *searcher) printStdout(blurb, link string) {
	if len(blurb) > 200 {
		blurb = blurb[:200]
	}
	if s.config.urls && len(blurb) > 0 {
		fmt.Println(link)
	}
	fmt.Println(blurb)
	fmt.Println()
}
