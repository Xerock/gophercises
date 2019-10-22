package link

import (
	"io"
	"strings"

	"golang.org/x/net/html"
)

// Link represent a html link and its content
type Link struct {
	Href, Text string
}

// Parse will take an HTML document and will return a slice of links parsed from it
func Parse(r io.Reader) (links []Link, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return
	}

	searchLinks(doc, &links)
	return links, nil
}

func searchLinks(n *html.Node, links *[]Link) {
	// If node is a link
	if n.Type == html.ElementNode && n.Data == "a" {
		// Find the href attribute
		for _, a := range n.Attr {
			if a.Key == "href" {
				text := strings.Join(strings.Fields(parseLinkContent(n)), " ")
				*links = append(*links, Link{a.Val, text})
				break
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		searchLinks(c, links)
	}
}

func parseLinkContent(n *html.Node) (content string) {
	for child := n.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.TextNode {
			content += child.Data
		}
		content += parseLinkContent(child)
	}
	return content
}
