package v2

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"golang.org/x/sync/errgroup"

	dwd "microservice/internal/dwd/v2/internal"
	"microservice/internal/dwd/v2/internal/parser"
)

var (
	errUnknownDatabase = errors.New("unknown database")
)

// DownloadFiles tries to download all available files for the given parameters.
// It returns the filepaths of the downloaded datafiles and (if availalbe) the
// description pages for the datasets.
func DownloadFiles(database, stationID string, product Product, granularity Granularity) (datafiles []string, descriptions [][2]string, err error) { //nolint:lll
	keys := make([]string, 0, len(Products))
	for k := range Products {
		keys = append(keys, k)
	}
	if !slices.Contains(keys, database) {
		return nil, nil, errUnknownDatabase
	}

	if !slices.Contains(Products[database][granularity], product) {
		return nil, nil, errUnsupportedProduct
	}

	uri, err := url.JoinPath(Databases[database], granularity.UrlPart(), product.UrlPart())
	if err != nil {
		return nil, nil, err
	}

	res, err := http.Get(uri) //nolint:gosec
	if err != nil {
		return nil, nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("%s: %w", uri, errStatusNotOK)
	}

	page, err := parser.ReadPage(res.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", uri, err)
	}

	possibleDescriptionFiles := parser.ParseFileLinks(page)
	possibleDataFolders := parser.ParseFolderLinks(page)

	descriptionFiles := make([][2]string, len(possibleDescriptionFiles))
	dataFiles := make([]string, 0)

	for idx, file := range possibleDescriptionFiles {
		uri, err := url.JoinPath(uri, file)
		if err != nil {
			return nil, nil, err
		}
		filepath, err := dwd.Download(uri)
		if err != nil {
			return nil, nil, err
		}

		descriptionFiles[idx] = [2]string{file, filepath}
	}

	var group errgroup.Group
	var l sync.Mutex

	for _, folder := range possibleDataFolders {
		group.Go(func() error {
			uri, err := url.JoinPath(uri, folder)
			if err != nil {
				return err
			}
			res, err := http.Get(uri) //nolint:gosec
			if err != nil {
				return err
			}
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("%s: %w", uri, errStatusNotOK)
			}

			page, err := parser.ReadPage(res.Body)
			if err != nil {
				return fmt.Errorf("%s: %w", uri, err)
			}

			possibleDataFiles := parser.ParseFileLinks(page)
			for _, datafile := range possibleDataFiles {
				if !strings.Contains(datafile, stationID) {
					continue
				}
				uri, err := url.JoinPath(uri, datafile)
				if err != nil {
					return err
				}
				filepath, err := dwd.Download(uri)
				if err != nil {
					return err
				}

				l.Lock()
				dataFiles = append(dataFiles, filepath)
				l.Unlock()
			}
			return nil

		})
	}

	if err := group.Wait(); err != nil {
		return nil, nil, err
	}

	return dataFiles, descriptionFiles, nil
}
