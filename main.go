package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Printf("You need to specify the url\n ./downloader [url]\n")
		os.Exit(3)
	}
	url := args[1]
	body := GetHtmlPage(url)
	fmt.Println(GetMediaFromHtml(body))
}

func GetHtmlPage(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode > 299 {
		log.Fatalf("Response failed with status code: %d and\nbody: %s\n", resp.StatusCode, body)
	}
	if err != nil {
		log.Fatal(err)
	} else {
		log.Printf("HTML body loaded, lenth: %v\n", len(body))
	}
	return body
}

func MediaTypeHandler(node *html.Node) string {
	switch {
	case node.Type == html.ElementNode && node.Data == "a":
		for _, a := range node.Attr {
			if a.Key == "href" {
				return a.Val
			}
		}
	case node.Type == html.ElementNode && node.Data == "img":
		for _, a := range node.Attr {
			if a.Key == "src" {
				return a.Val
			}
		}
	}
	return ""
}

func GetMediaFromHtml(body []byte) []string {
	var links []string

	doc, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	var f func(*html.Node)
	f = func(n *html.Node) {
		url := MediaTypeHandler(n)
		if len(url) > 0 {
			links = append(links, url)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	return links
}
