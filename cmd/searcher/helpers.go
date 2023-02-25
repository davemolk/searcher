package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

// getTerms reads terms off stdin and sets the value
// of s.terms to the resulting slice.
func (s *searcher) getTerms() {
	var terms []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// handle phrases
		term := strings.Replace(scanner.Text(), " ", "+", -1)
		terms = append(terms, term)
	}
	if scanner.Err() != nil {
		s.errorLog.Fatal("problem reading terms off stdin: ", scanner.Err())
		return
	}
	s.terms = terms
}

// encode takes in a map, encodes it to json, and returns
// a byte slice and any errors.
func (s *searcher) encode(data map[string]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	return bytes.TrimRight(buf.Bytes(), "\n"), err
}

// processResults encodes the search results and either prints
// them as json to stdout or writes them to a file (or both).
func (s *searcher) processResults() {
	switch {
	case len(s.searches.search) > 0:
		b, err := s.encode(s.searches.search)
		if err != nil {
			s.errorLog.Printf("unable to encode search map to json: %v\n", err)
			return
		}
		if s.config.json {
			fmt.Println(string(b))
		}
		if s.config.write {
			if err := os.WriteFile("data/search.json", b, 0644); err != nil {
				s.errorLog.Printf("write error: %v\n", err)
			}
		}
	case len(s.searches.searches) > 0:
		var wg sync.WaitGroup
		results := make(chan string, len(s.terms))
		for _, t := range s.terms {
			wg.Add(1)
			go func(t string) {
				defer wg.Done()
				b, err := s.encode(s.searches.searches[t])
				if err != nil {
					s.errorLog.Printf("unable to encode searches map to json: %v\n", err)
					return
				}
				results <- string(b)
				if s.config.write {
					name := fmt.Sprintf("data/%s.json", t)
					if err := os.WriteFile(name, b, 0644); err != nil {
						s.errorLog.Printf("write error: %v\n", err)
					}
				}
			}(t)

		}
		wg.Wait()
		close(results)
		if s.config.json {
			var j []string
			for r := range results {
				j = append(j, r)
			}
			fmt.Println(j)
		}
	}
}
