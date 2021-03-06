package main

import (
	"flag"
	"fmt"

	"github.com/gocolly/colly"
)

type Node struct {
	parent   *Node
	data     string
	children []*Node
}

func (n *Node) VisitChildren(c *colly.Collector) {
	for index := range n.children {
		c.Visit(n.children[index].data)
	}
}

func Iterate(start Node) chan Node {
	ch := make(chan Node)

	ch <- start

	go func(ch chan Node) {
		for index := range start.children {
			ch <- *start.children[index]
		}

		close(ch)
	}(ch)

	return ch
}

// main() contains code adapted from example found in Colly's docs:
// http://go-colly.org/docs/examples/basic/
func main() {
	var startNode Node
	var startLink string
	wikiFlagDescription := "should be followed by a wikipedia page to start from"
	flag.StringVar(&startLink, "wiki", "hello", wikiFlagDescription)
	flag.Parse()

	startNode.data = startLink

	// Instantiate default collector
	c := colly.NewCollector()

	// On every a element which has href attribute call callback
	c.OnHTML(`#content a[href^="/wiki"]`, func(e *colly.HTMLElement) {
		link := e.Attr("href")

		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		var newNode Node
		newNode.parent = &startNode
		newNode.data = "https://en.wikipedia.org" + link
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	finalNode, err := loop(c, startNode)

	if err == nil {
		var path []string
		for currentNode := finalNode; currentNode != nil; currentNode = currentNode.parent {
			path = append([]string{currentNode.data}, path...)
		}

		for index := range path {
			fmt.Println(path[index])
		}
	} else {
		fmt.Println("No Bacon!")
	}
}

func loop(c *colly.Collector, start Node) (*Node, error) {
	ch := make(chan Node)

	ch <- start

	for range ch {
		for node := range Iterate(start) {
			if node.data == "https://en.wikipedia.org/wiki/Kevin_Bacon" {
				return &node, nil
			}
			c.Visit(node.data)
			for index := range node.children {
				ch <- *node.children[index]
			}
		}
	}
	return &Node{}, nil
}
