package helpers

import (
	"strings"

	"golang.org/x/net/html"
)

// FilterDocumentForFolders iterates recursively through the document nodes and
// tries to find link elements (<a>) which end with a slash ("/").
// This function excludes the folder path "../".
func FilterDocumentForFolders(document *html.Node) (folders []string) {
	var filter func(node *html.Node)
	filter = func(node *html.Node) {
		// check if the current node is an element and is an <a> element
		if node.Type == html.ElementNode && node.Data == "a" {
			// since the node is a link to a possible folder, get the link to it
			link := node.FirstChild.Data
			if link != "../" && strings.HasSuffix(link, "/") {
				folders = append(folders, link)
			}
		}
		// now iterate through other possible nodes of the element and filter
		// though them as wenn with this filter
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			filter(child)
		}
	}
	filter(document)
	return folders
}

// FilterDocumentForFiles iterates recursively through the document nodes and
// tries to find link elements in it which are files (i.e., not ending with a
// slash)
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
