package main

import (
	"fmt"
	"strings"
	"sync"
)

type queryData struct {
	base   string
	spacer string
}

type parseData struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
}

func (s *searcher) makeQueryData() []queryData {
	if s.config.verbose {
		s.infoLog.Println("Making slice of query data...")
	}
	var qdSlice []queryData

	ask := queryData{
		base:   "https://www.ask.com/web?q=",
		spacer: "+",
	}

	bing := queryData{
		base:   "https://bing.com/search?q=",
		spacer: "+",
	}

	brave := queryData{
		base:   "https://search.brave.com/search?q=",
		spacer: "+",
	}

	duck := queryData{
		base:   "https://html.duckduckgo.com/html?q=",
		spacer: "+",
	}

	yahoo := queryData{
		base:   "https://search.yahoo.com/search?p=",
		spacer: "+",
	}

	yandex := queryData{
		base:   "https://yandex.com/search/?text=",
		spacer: "+",
	}
	qdSlice = append(qdSlice, ask, bing, brave, duck, yahoo, yandex)

	return qdSlice
}

func (s *searcher) cleanQuery() {
	// handle multiple words
	s.config.baseSearch = strings.Replace(s.config.baseSearch, " ", "+", -1)
	if s.config.exact {
		s.config.baseSearch = fmt.Sprintf("\"%s\"", s.config.baseSearch)
	}
}

func (s *searcher) makeSearchURLs() [6]chan string {
	// each channel will store all the query strings for a given search engine.
	var chans [6]chan string
	for i := range chans {
		chans[i] = make(chan string, len(s.terms))
	}

	qdSlice := s.makeQueryData()

	var wg sync.WaitGroup
	for _, term := range s.terms {
		for i, qd := range qdSlice {
			wg.Add(1)
			go func(qd queryData, term string, i int) {
				defer wg.Done()
				url := fmt.Sprintf("%s%s%s%s", qd.base, s.config.baseSearch, qd.spacer, term)
				// jenky, lol
				url = fmt.Sprintf("%sGETTERM%s", url, term)
				chans[i] <- url
			}(qd, term, i)
		}
	}

	wg.Wait()
	for i := range chans {
		close(chans[i])
	}

	if s.config.verbose {
		s.infoLog.Println("Search URLs created.")
	}

	return chans
}
