package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func parseLinks(document *html.Node) (links []string) {
	var filter func(node *html.Node)
	filter = func(node *html.Node) {
		if node.Type != html.ElementNode || node.Data != "a" {
			goto filterAgain
		}

		for _, attr := range node.Attr {
			if attr.Key != "href" {
				continue
			}

			links = append(links, strings.TrimSpace(attr.Val))

		}

	filterAgain:
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			filter(child)
		}
	}
	filter(document)
	return
}
