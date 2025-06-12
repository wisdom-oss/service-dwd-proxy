package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseFileLinks(document *html.Node) (relativeUrls []string) {
	links := parseLinks(document)

	for _, link := range links {
		if link == "" || link == "../" || strings.HasSuffix(link, "/") {
			continue
		}

		relativeUrls = append(relativeUrls, strings.TrimSpace(link))
	}
	return
}
