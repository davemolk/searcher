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

	yandex := queryData{
		base:   fmt.Sprintf("%s%s", "https://yandex.com/search/?text=", s.config.baseSearch),
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
