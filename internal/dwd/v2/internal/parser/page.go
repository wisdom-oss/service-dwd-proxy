package parser

import (
	"io"

	"golang.org/x/net/html"
)

func ReadPage(r io.Reader) (*html.Node, error) {
	return html.ParseWithOptions(r, html.ParseOptionEnableScripting(false))
}
