package main

import (
	"fmt"
	"time"
)

//step 3:define the interface to obtain web data from any type of data
type Fetcher interface {
	//pass URL to Fetch and return web body,urls in the page
	//represented by this URL,and a err
	Fetch(url string) (body string, urls []string, err error)
}

//step 6:pass the web url,the recursion depth of find web,the interface of get
//web content.
//Crawl to get web content
//         deal with web content
//         display get,deal result

//Fetch(url)->var->channels->another var
func Crawl(url string, depth int, fetcher Fetcher) {
	type URL struct { //store the current url's recursion depth count
		url   string
		depth int
	}

	msg := make(chan string) //display result channel
	req := make(chan URL)    //current url's recursion depth count channel
	quit := make(chan int)   //?

	//define crawler in this way to get web content,
	//easy to use routines and routines channels
	crawler := func(url string, depth int) {
		defer func() { quit <- 0 }() //run at the end of crawler

		if depth <= 0 { //when depth=0,we not recursion
			return
		}

		body, urls, err := fetcher.Fetch(url) //get web content

		if err != nil {
			msg <- fmt.Sprintf("%s\n", err) //fill display error result channel
			return
		}

		msg <- fmt.Sprintf("found: %s %q\n", url, body) //fill display result channel

		//fill recursion depth count struct channel
		for _, u := range urls {
			req <- URL{u, depth - 1} //depth count down
		}
	}

	works := 1

	memo := make(map[string]bool)
	memo[url] = true

	//run a crawler routine to fill channels
	go crawler(url, depth)

	for works > 0 {
		//check and select channels to recursion
		select {
		//output "get web content info"
		case s := <-msg:
			fmt.Print(s)
		//recursion to get web content
		case u := <-req:
			if !memo[u.url] {
				memo[u.url] = true
				works++

				go crawler(u.url, u.depth)
			}
		//when all crawler routine finish the for is finish
		case <-quit:
			works--
		}
	}
}

func main() {
	//pass parameter:fetcher Fetcher=fetcher
	//it's the use of interface
	Crawl("http://golang.org/", 4, fetcher)
}

//step 2:map data type of "string:fakeResult"
type fakeFetcher map[string]*fakeResult

//step 1:data type to store a web page content
//contain the web body and the urls in this page
type fakeResult struct {
	body string
	urls []string
}

//step 5:the method of fakeResult that implement interface Fetcher.Fetch
//the function of this method must same as Fetcher.Fetch
func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(time.Second * 5)
	if res, ok := f[url]; ok { //define a var in this format can use ok to                                     //decide if the url key have value or if the key exist
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

//Step 4:filling the data type that store web content.
//In this sample we manual fill the body and urls and all  recursion to simulation a web
//define the var of fakeFetcher type is fetcher
//fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"http://golang.org/": &fakeResult{
		"The Go Programming Language", //fakeResult body
		[]string{ //fakeResult urls
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
