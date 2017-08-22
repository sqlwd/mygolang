package main

import (
	"fmt"
	"time"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

func Crawl(url string, depth int, fetcher Fetcher) {
	type URL struct {
		url   string
		depth int
	}

	msg := make(chan string)
	req := make(chan URL)
	quit := make(chan int)

	crawler := func(url string, depth int) {
		defer func() { quit <- 0 }()

		if depth <= 0 {
			return
		}

		body, urls, err := fetcher.Fetch(url)

		if err != nil {
			msg <- fmt.Sprintf("%s\n", err)
			return
		}

		msg <- fmt.Sprintf("found: %s %q\n", url, body)

		for _, u := range urls {
			req <- URL{u, depth - 1}
		}
	}

	works := 1

	memo := make(map[string]bool)
	memo[url] = true

	go crawler(url, depth)

	for works > 0 {
		select {
		case s := <-msg:
			fmt.Print(s)
		case u := <-req:
			if !memo[u.url] {
				memo[u.url] = true
				works++

				go crawler(u.url, u.depth)
			}
		case <-quit:
			works--
		}
	}
}

func main() {
	Crawl("http://golang.org/", 4, fetcher)
}

type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(time.Second * 5)
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"http://golang.org/pkg/",
			"http://golang.org/cmd/",
		},
	},
	"http://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"http://golang.org/",
			"http://golang.org/cmd/",
			"http://golang.org/pkg/fmt/",
			"http://golang.org/pkg/os/",
		},
	},
	"http://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
	"http://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"http://golang.org/",
			"http://golang.org/pkg/",
		},
	},
}
