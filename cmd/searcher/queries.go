package main

import (
	"fmt"
	"strings"
	"sync"
)

// queryData is a struct containing the base URL
// and spacer to be used in constructing query strings.
type queryData struct {
	base   string
	spacer string
}

// parseData is a struct containing information on
// parsing search engine results for each URL and blurb.
type parseData struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
}

// makeQueryData returns a slice of queryData featuring
// one instance for each search engine (ask, bing, brave,
// duck duck go, and yahoo).
func (s *searcher) makeQueryData() []queryData {
	var qdSlice []queryData

	ask := queryData{
		base:   fmt.Sprintf("%s%s", "https://www.ask.com/web?q=", s.config.baseSearch),
		spacer: "+",
	}

	bing := queryData{
		base:   fmt.Sprintf("%s%s", "https://bing.com/search?q=", s.config.baseSearch),
		spacer: "+",
	}

	brave := queryData{
		base:   fmt.Sprintf("%s%s", "https://search.brave.com/search?q=", s.config.baseSearch),
		spacer: "+",
	}

	duck := queryData{
		base:   fmt.Sprintf("%s%s", "https://html.duckduckgo.com/html?q=", s.config.baseSearch),
		spacer: "+",
	}

	yahoo := queryData{
		base:   fmt.Sprintf("%s%s", "https://search.yahoo.com/search?p=", s.config.baseSearch),
		spacer: "+",
	}

	qdSlice = append(qdSlice, ask, bing, brave, duck, yahoo)

	return qdSlice
}

// cleanQuery replaces any spaces with "+" and adds double quotes when
// the exact flag is invoked on the command line.
func (s *searcher) cleanQuery() {
	// handle multiple words
	s.config.baseSearch = strings.Replace(s.config.baseSearch, " ", "+", -1)
	if s.config.exact {
		s.config.baseSearch = fmt.Sprintf("\"%s\"", s.config.baseSearch)
	}
}

// makeSearchURLs returns an array of 5 string channels, each containing all
// the query strings for a given search engine.
func (s *searcher) makeSearchURLs() [5]chan string {
	var chans [5]chan string
	qdSlice := s.makeQueryData()
	var wg sync.WaitGroup

	switch {
	case len(s.terms) == 0:
		for i := range chans {
			chans[i] = make(chan string, 1)
		}
		for i, qd := range qdSlice {
			wg.Add(1)
			go func(qd queryData, i int) {
				defer wg.Done()
				chans[i] <- qd.base
			}(qd, i)
		}
	default:
		for i := range chans {
			chans[i] = make(chan string, len(s.terms))
		}
		for _, term := range s.terms {
			for i, qd := range qdSlice {
				wg.Add(1)
				go func(qd queryData, term string, i int) {
					defer wg.Done()
					url := fmt.Sprintf("%s%s%s", qd.base, qd.spacer, term)
					// jenky, lol
					url = fmt.Sprintf("%sGETTERM%s", url, term)
					chans[i] <- url
				}(qd, term, i)
			}
		}
	}

	wg.Wait()
	for i := range chans {
		close(chans[i])
	}

	return chans
}
