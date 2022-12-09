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

func (s *searcher) getTerms() {
	switch {
	case s.config.terms:
		terms, err := s.readInput()
		if err != nil {
			s.errorLog.Fatalf("unable to read terms off stdin: %v", err)
		}
		s.terms = terms
	default:
		if s.config.verbose {
			s.errorLog.Println("No additional search terms supplied. Continuing with base search only.")
		}
	}
}

// readInput reads terms off stdin and returns
// a string slice containing those terms.
func (s *searcher) readInput() ([]string, error) {
	var terms []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		// handle phrases
		term := strings.Replace(scanner.Text(), " ", "+", -1)
		terms = append(terms, term)
	}
	return terms, scanner.Err()
}

func (s *searcher) encode(data map[string]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	return bytes.TrimRight(buf.Bytes(), "\n"), err
}

func (s *searcher) dump() {
	switch {
	case len(s.searches.search) > 0:
		b, err := s.encode(s.searches.search)
		if err != nil {
			s.errorLog.Printf("unable to encode map to json: %v\n", err)
			return
		}
		if s.config.json {
			fmt.Println(string(b))
		}
		if s.config.write {
			err := os.WriteFile("data/search.json", b, 0644)
			if err != nil {
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
					s.errorLog.Println("unable to encode map to json", err)
					return
				}
				results <- string(b)
				if s.config.write {
					name := fmt.Sprintf("data/%s.json", t)
					err := os.WriteFile(name, b, 0644)
					if err != nil {
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
