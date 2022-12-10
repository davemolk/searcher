package main

import "sync"

// searchMap contains two maps, search and searches. The first is
// used when a single query is entered (format URL:blurb), while
// the second is used in situations where the -terms flag is true
// (format term:URL:blurb).
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

// storeSearches checks if a URL has already been stored for
// a given search term. If it hasn't, storeSearches will 
// store the URL and blurb and associate them with the search term.
func (s *searchMap) storeSearches(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; !ok {
		s.searches[term][url] = blurb
	}
}

// storeSearch accepts a URL and associated blurb, checks whether
// or not the URL is already present in the map. and stores the
// data if not.
func (s *searchMap) storeSearch(url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.search[url]; !ok {
		s.search[url] = blurb
	}
}
