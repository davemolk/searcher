package main

import (
	"flag"
	"log"
	"os"
	"regexp"
)

type config struct {
	baseSearch  string
	concurrency int
	exact       bool
	file        string
	timeout     int
	write       bool
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
	flag.BoolVar(&config.exact, "se", false, "search for exact match")
	flag.StringVar(&config.baseSearch, "s", "", "base search (enclose phrases in quotes)")
	flag.StringVar(&config.file, "f", "", "file name containing additional terms to run with the base search")
	flag.IntVar(&config.timeout, "t", 5000, "timeout (in ms, default 5000)")
	flag.BoolVar(&config.write, "w", false, "write results to file")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)

	searches := newSearchMap()
	noBlank := regexp.MustCompile(`\s{2,}`)

	if config.write {
		err := os.Mkdir("data", 0755)
		if err != nil {
			log.Fatal(err)
		}
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

	if config.write {
		wait := s.launchWriters()
		<-wait
	}
}
