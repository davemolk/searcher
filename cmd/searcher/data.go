package main

import "sync"

// searchMap contains two maps, search and searches. The first is
// used when a single query is entered (map formatted as URL:blurb),
// while the second is used when the -terms flag is true
// (formatted as a map of maps, with term:URL:blurb).
type searchMap struct {
	mu       sync.Mutex
	search   map[string]string
	searches map[string]map[string]string
}

// newSearchMap initializes and returns a new searchMap.
func newSearchMap() *searchMap {
	return &searchMap{
		searches: make(map[string]map[string]string),
		search:   make(map[string]string),
	}
}

// Given search term, storeSearches checks if a URL has already
// been stored. If it hasn't, the URL and blurb will be
// stored as a map value for the term.
func (s *searchMap) storeSearches(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; !ok {
		s.searches[term][url] = blurb
	}
}

// storeSearch checks if a URL has already been stored. If it hasn't,
// both the URL and the associated blurb will be stored.
func (s *searchMap) storeSearch(url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.search[url]; !ok {
		s.search[url] = blurb
	}
}
