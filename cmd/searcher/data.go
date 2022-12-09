package main

import "sync"

// searchMap contains two maps. search is used when only a
// single query is entered, while searches is used when
// the -terms flag is true.
type searchMap struct {
	mu       sync.Mutex
	search   map[string]string
	searches map[string]map[string]string
}

func newSearchMap() *searchMap {
	return &searchMap{
		searches: make(map[string]map[string]string),
		search:   make(map[string]string),
	}
}

func (s *searchMap) store(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; !ok {
		s.searches[term][url] = blurb
	}
}

func (s *searchMap) storeSearch(url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.search[url]; !ok {
		s.search[url] = blurb
	}
}
