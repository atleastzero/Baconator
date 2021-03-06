package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
)

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	var currentLink string
	wikiFlagDescription := "should be followed by a wikipedia page to start from"
	flag.StringVar(&currentLink, "wiki", "hello", wikiFlagDescription)
	flag.Parse()
	// degreeArray := make(map[string]string)

	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML(`#content a[href^="/wiki"]`, func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit("https://en.wikipedia.org" + link)
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on currentLink
	c.Visit(currentLink)
}
