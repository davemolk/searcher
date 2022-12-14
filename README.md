# searcher

Run a base query (plus optional add-ons) through ask, bing, brave, duck duck go, yahoo, and yandex.

## Overview
By default, this tool collects the URLs and result blurbs and prints them to stdout. You can pipe in additional terms that will be added to the base query. Print the results as json if you'd like, or save them to json files.

## Examples
```
$ cat terms.txt
microservices
mascot
cloud
cli
```
Let's use *golang* as our main search query and combine it with each of the above terms. We'll encode the search results as json and print to stdout and also save them as json files.

$ cat terms.txt | searcher -q golang -j -w -t
(make sure you include -t or the terms.txt won't be picked up!)
```
https://search.brave.com/search?q=golang+mascot
https://bing.com/search?q=golang+cloud
https://bing.com/search?q=golang+mascot
https://www.ask.com/web?q=golang+microservices
https://www.ask.com/web?q=golang+mascot
https://www.ask.com/web?q=golang+cli
https://search.brave.com/search?q=golang+cli
https://www.ask.com/web?q=golang+cloud
https://bing.com/search?q=golang+cli
https://bing.com/search?q=golang+microservices
https://search.brave.com/search?q=golang+microservices
https://search.brave.com/search?q=golang+cloud
https://html.duckduckgo.com/html?q=golang+cloud
https://html.duckduckgo.com/html?q=golang+cli
https://html.duckduckgo.com/html?q=golang+mascot
https://html.duckduckgo.com/html?q=golang+microservices
https://search.yahoo.com/search?p=golang+cloud
https://search.yahoo.com/search?p=golang+cli
https://search.yahoo.com/search?p=golang+mascot
https://search.yahoo.com/search?p=golang+microservices
https://yandex.com/search/?text=golang+cli
https://yandex.com/search/?text=golang+cloud
https://yandex.com/search/?text=golang+mascot
https://yandex.com/search/?text=golang+microservices
```
the results are printed as JSON to the stdout and saved in the following files:
```
cli.json
cloud.json
mascot.json
microservices.json
```
where each file contains the JSON object (URL:blurb) for that particular term.

## Install
First, you'll need to [install go](https://golang.org/doc/install). Then, run the following command:

```
go install github.com/davemolk/searcher/cmd/searcher@latest
```

## Flags
```
Usage of searcher:
  -c int
    	Max number of goroutines to use at any given time.
  -e bool
    	Exact matching for query.
  -j bool
    	Print results as JSON.
  -q string
    	Base search query.
  -t bool
    	Check Stdin for additional search terms.
  -to int
    	Request timeout (in ms).
  -v bool
    	Verbose output.
  -w bool
    	Write results to file.
```
