package main

import "sync"

// searchesMap has the search term(s) as the key(s) and a
// nested map as the value(s). The nested map is in the
// form URL: blurb.
type searchesMap struct {
	mu       sync.Mutex
	searches map[string]map[string]string
}

func newSearchMap() *searchesMap {
	return &searchesMap{
		searches: make(map[string]map[string]string),
	}
}

func (s *searchesMap) store(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; ok {
		return
	}
	s.searches[term][url] = blurb
}