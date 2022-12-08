package main

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// makeParseData returns the selectors for each of the search engines.
func (s *searcher) makeParseData() []parseData {
	if s.config.verbose {
		s.infoLog.Println("Making slice of parse data...")
	}
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

	yandex := parseData{
		blurbSelector: "div.TextContainer",
		itemSelector:  "li.serp-item",
		linkSelector:  "div.VanillaReact > a",
		name:          "yandex",
	}

	pdSlice = append(pdSlice, ask, bing, brave, duck, yahoo, yandex)

	return pdSlice
}

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
				body, err := s.makeRequest(urlTerm[0], s.config.timeout)
				if err != nil {
					if s.config.verbose {
						s.errorLog.Printf("error in makeRequest: %v\n", err)
					}
					<-tokens
					return
				}
				<-tokens
				s.parseSearchResults(body, urlTerm[1], pdSlice[i])
			}(u, i)
		}
	}

	wg.Wait()
}

func (s *searcher) parseSearchResults(data, term string, pd parseData) {
	if s.config.verbose {
		s.infoLog.Printf("Parsing %s for %q", pd.name, term)
	}
	localResults := make(map[string]string)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
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
		if blurb == "" && s.config.verbose {
			s.errorLog.Printf("unable to get blurb for %s\n", pd.name)
		}
		cleanedLink := s.cleanLinks(link)
		cleanedBlurb := s.cleanBlurb(blurb)
		fmt.Println(cleanedBlurb)
		fmt.Println()
		localResults[cleanedLink] = cleanedBlurb
		s.searches.store(term, cleanedLink, cleanedBlurb)
	})
}

func (s *searcher) cleanBlurb(str string) string {
	cleanB := s.noBlank.ReplaceAllString(str, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

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
