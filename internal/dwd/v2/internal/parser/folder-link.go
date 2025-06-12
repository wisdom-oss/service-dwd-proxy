package parser

import (
	"strings"

	"golang.org/x/net/html"
)

func ParseFolderLinks(document *html.Node) (relativeFolders []string) {
	links := parseLinks(document)

	for _, link := range links {
		if link == "" || !strings.HasSuffix(link, "/") {
			continue
		}

		relativeFolders = append(relativeFolders, strings.TrimSpace(link))
	}

	return
}
