package dwd

import (
	"strings"

	"golang.org/x/net/html"
)

func FilterDocumentForFiles(document *html.Node) (files []string) {
	var filter func(node *html.Node)
	filter = func(node *html.Node) {
		// check if the current node is an element and is an <a> element
		if node.Type == html.ElementNode && node.Data == "a" {
			// since the node is a link to a possible folder, get the link to it
			var link string
			for _, attr := range node.Attr {
				if attr.Key == "href" {
					link = attr.Val
				}
			}
			if link != "" && link != "../" && !strings.HasSuffix(link, "/") {
				files = append(files, link)
			}
		}
		// now iterate through other possible nodes of the element and filter
		// though them as wenn with this filter
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			filter(child)
		}
	}
	filter(document)
	return files
}
