package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
)

// readInput takes in the file name for a list of terms and returns
// a string slice containing those terms.
func (s *searcher) readInput(name string) ([]string, error) {
	var terms []string
	n, err := os.Open(name)
	if err != nil {
		return terms, err
	}
	defer n.Close()

	scanner := bufio.NewScanner(n)
	for scanner.Scan() {
		terms = append(terms, scanner.Text())
	}
	return terms, scanner.Err()
}

// getTerms looks at the user flag input, determines whether a single
// term or a file name for a list of terms has been selected, and
// adds the appropriate field to the fof struct instance.
func (s *searcher) getTerms() {
	switch {
	case s.config.file != "":
		terms, err := s.readInput(s.config.file)
		if err != nil {
			s.errorLog.Fatalf("unable to get terms from file: %v", err)
		}
		s.terms = terms
	default:
		s.errorLog.Println("No search terms supplied. Continuing with search target only.")
	}
}

func (s *searcher) writeData(name string, data map[string]string) {
	file, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	b, err := s.encode(data)
	if err != nil {
		log.Fatal(err)
	}
	_, err = file.Write(b)
	if err != nil {
		log.Fatal(err)
	}
	err = file.Sync()
	if err != nil {
		log.Fatal(err)
	}
}

func (s *searcher) encode(data map[string]string) ([]byte, error) {
	buf := &bytes.Buffer{}
	encoder := json.NewEncoder(buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(data)
	return bytes.TrimRight(buf.Bytes(), "\n"), err
}
