package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	link "gophercises/html-link-parser"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	var urlFlag *string = flag.String("url", "https://gophercises.com", "the website url to crawl")
	var maxDepth *int = flag.Int("depth", 10, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	toXML := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXML.Urls = append(toXML.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXML); err != nil {
		panic(err)
	}
	fmt.Println()
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

func get(urlStr string) []string {
	//fmt.Println("Getting", urlStr)
	resp, err := http.Get(urlStr)
	if err != nil {
		fmt.Println("Failed to get", urlStr)
		os.Exit(1)
	}
	defer resp.Body.Close()

	reqURL := resp.Request.URL
	baseURL := &url.URL{
		Scheme: reqURL.Scheme,
		Host:   reqURL.Host,
	}
	base := baseURL.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func bfs(urlStr string, maxDepth int) []string {
	visited := make(map[string]bool)
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}

	for i := 0; i <= maxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, present := visited[url]; !present {
				visited[url] = true
				for _, link := range get(url) {
					nq[link] = struct{}{}
				}
			}
		}
	}

	ret := make([]string, 0, len(visited))
	for url := range visited {
		ret = append(ret, url)
	}
	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}
	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}

// func buildMap(s string) (sitemap string) {
// 	// TODO : Handle when url does or doesn't end with '/'
// 	visited[s] = true
// 	crawl(s)

// 	sitemap += `<?xml version="1.0" encoding="UTF-8"?>
// <urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
// `

// 	for site := range visited {
// 		output, err := xml.MarshalIndent(url{site}, "  ", "    ")
// 		if err != nil {
// 			fmt.Printf("error: %v\n", err)
// 		}
// 		sitemap += string(output) + "\n"
// 	}
// 	sitemap += "</urlset>"
// 	return sitemap
// }
