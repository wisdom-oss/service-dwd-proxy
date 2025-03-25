package dwd

import (
	"fmt"
	"strings"

	"golang.org/x/net/html"
)

func GetFolderURLs(p *html.Node, baseUrl string) (urls []string) {
	var filter func(node *html.Node)
	filter = func(node *html.Node) {
		// check if the current node is an element and is an <a> element
		if node.Type == html.ElementNode && node.Data == "a" {
			// since the node is a link to a possible folder, get the link to it
			link := node.FirstChild.Data
			if link != "../" && strings.HasSuffix(link, "/") {
				urls = append(urls, fmt.Sprintf(`%s/%s`, baseUrl, link))
			}
		}
		// now iterate through other possible nodes of the element and filter
		// though them as wenn with this filter
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			filter(child)
		}

	}
	filter(p)
	return urls
}
