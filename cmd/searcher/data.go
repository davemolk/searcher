package main

import "sync"

// searchesMap has the search term(s) as the key(s) and a
// nested map as the value(s). The nested map is in the
// form URL: blurb.
type searchMap struct {
	mu       sync.Mutex
	searches map[string]map[string]string
}

func newSearchMap() *searchMap {
	return &searchMap{
		searches: make(map[string]map[string]string),
	}
}

func (s *searchMap) store(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; ok {
		return
	}
	s.searches[term][url] = blurb
}
