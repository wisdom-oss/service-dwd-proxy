package dwd

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

var ErrNotFound = errors.New("page not found")
var ErrResponseNotOK = errors.New("response not 200")

func LoadIndexPage(url string) (*html.Node, error) {
	res, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("index page request failed: %w", err)
	}

	if res.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if res.StatusCode != http.StatusOK {
		return nil, ErrResponseNotOK
	}

	document, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	return document, nil
}
