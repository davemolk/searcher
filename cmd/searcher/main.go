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
	terms       bool
	timeout     int
	verbose     bool
	write       bool
}

type searcher struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
	noBlank  *regexp.Regexp
	searches *searchMap
	terms    []string
}

func main() {
	var config config
	flag.StringVar(&config.baseSearch, "q", "", "base search query")
	flag.IntVar(&config.concurrency, "c", 10, "max number of goroutines to use at any given time")
	flag.BoolVar(&config.exact, "e", false, "search for exact match")
	flag.BoolVar(&config.terms, "terms", false, "check stdin for additional search terms")
	flag.IntVar(&config.timeout, "t", 5000, "timeout (in ms, default 5000)")
	flag.BoolVar(&config.verbose, "v", false, "verbose output")
	flag.BoolVar(&config.write, "w", false, "write results to file")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)

	searches := newSearchMap()
	noBlank := regexp.MustCompile(`\s{2,}`)

	if config.write {
		if err := os.Mkdir("data", 0755); err != nil {
			log.Fatal("can't make data folder", err)
		}
	}

	s := &searcher{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		noBlank:  noBlank,
		searches: searches,
	}

	s.cleanQuery()
	s.getTerms()

	if config.terms {
		for _, t := range s.terms {
			s.searches.searches[t] = make(map[string]string)
		}
	}

	s.getAndParseData()

	if config.write {
		wait := s.launchWriters()
		<-wait
	}
}
