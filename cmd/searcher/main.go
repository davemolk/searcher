package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"sync"
)

type config struct {
	concurrency  int
	exact        bool
	searchTarget string
	file         string
	timeout      int
}

type searcher struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
	noBlank  *regexp.Regexp
	searches *searchesMap
	terms    []string
}

func main() {
	var config config
	flag.IntVar(&config.concurrency, "c", 10, "max number of goroutines to use at any given time")
	flag.BoolVar(&config.exact, "exact", false, "exact match of search query (some engines will only provide exact matches, while others will give 'close to exact' as well as exact)")
	flag.StringVar(&config.searchTarget, "s", "", "base search target (please enclose phrases in quotes)")
	flag.StringVar(&config.file, "f", "", "file name containing a list of terms")
	flag.IntVar(&config.timeout, "t", 5000, "timeout (in ms, default 5000)")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)

	searches := newSearchMap()
	noBlank := regexp.MustCompile(`\s{2,}`)

	err := os.Mkdir("data", 0755)
	if err != nil {
		log.Fatal(err)
	}

	s := &searcher{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		noBlank:  noBlank,
		searches: searches,
	}

	s.getTerms()
	for _, t := range s.terms {
		searches.searches[t] = make(map[string]string)
	}

	qdSlice := s.makeQueryData()
	pdSlice := s.makeParseData()

	chans := s.makeSearchURLs(qdSlice)

	s.getAndParseData(pdSlice, chans)

	var wg sync.WaitGroup
	for _, t := range s.terms {
		wg.Add(1)
		go func(t string) {
			defer wg.Done()
			name := fmt.Sprintf("data/%s.json", t)
			s.writeData(name, s.searches.searches[t])
		}(t)
	}

	wg.Wait()
}
