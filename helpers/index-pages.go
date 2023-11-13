package helpers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/html"
)

var ErrStatusNot200 = errors.New("response returned with a non-200 code")

func GetIndexPage(url string) (*html.Node, error) {
	// now request the index page generated by the open data portal
	log.Debug().Str("url", url).Msg("requesting index page")
	res, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to request index page '%s': %w", url, err)
	}
	// now check if the correct status code (200) has been returned
	if res.StatusCode != 200 {
		return nil, ErrStatusNot200
	}
	document, err := html.Parse(res.Body)
	if err != nil {
		log.Error().Err(err).Msg("unable to parse index page")
		return nil, fmt.Errorf("unable to parse index page '%s': %w", url, err)
	}
	return document, nil
}