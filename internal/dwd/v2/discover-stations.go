package v2

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"microservice/internal/dwd/v2/internal/parser"
	v2 "microservice/types/v2"
)

var (
	errUnsupportedProduct = errors.New("unsupported product in granularity")
	errStatusNotOK        = errors.New("the response code indicated a unsuccessful request")
)

func DiscoverStations(databaseUrl string, granularity Granularity, product Product) ([]v2.Station, error) {
	if !slices.Contains(AvailableClimateObservationProducts[granularity], product) {
		return nil, errUnsupportedProduct
	}

	uri, err := url.JoinPath(databaseUrl, granularity.UrlPart(), product.UrlPart())
	if err != nil {
		return nil, err
	}
	res, err := http.Get(uri) //nolint:gosec
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s: %w", uri, errStatusNotOK)
	}

	page, err := parser.ReadPage(res.Body)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", uri, err)
	}

	folders := parser.ParseFolderLinks(page)

	var eGroup errgroup.Group
	var stationFiles []string
	var arrayLock sync.Mutex

	for _, folder := range folders {
		if folder == "../" {
			continue
		}
		eGroup.Go(func() error {
			folderUrl, err := url.JoinPath(uri, folder)
			if err != nil {
				return err
			}
			res, err := http.Get(folderUrl) //nolint:gosec
			if err != nil {
				return err
			}
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("%s: %w", folderUrl, errStatusNotOK)
			}

			page, err := parser.ReadPage(res.Body)
			if err != nil {
				return fmt.Errorf("%s: %w", folderUrl, err)
			}

			files := parser.ParseFileLinks(page)
			for _, file := range files {
				if !strings.HasSuffix(file, "_Stationen.txt") {
					continue
				}

				{
					uri, err := url.JoinPath(folderUrl, file)
					if err != nil {
						return err
					}
					arrayLock.Lock()

					stationFiles = append(stationFiles, uri)
					arrayLock.Unlock()
				}

			}
			return nil
		})
	}

	if err := eGroup.Wait(); err != nil {
		return nil, err
	}

	var stations []v2.Station

	for _, file := range stationFiles {
		eGroup.Go(func() error {
			res, err := http.Get(file) //nolint:gosec
			if err != nil {
				return err
			}

			parsedStations, dateAreas, err := parser.ParseStationList(res.Body)
			if err != nil {
				return err
			}

			for i := range parsedStations {
				station := parsedStations[i]
				dateAvailabliliy := dateAreas[i]

				station.SupportedProducts = map[string]map[string][2]time.Time{
					product.String(): {
						granularity.String(): dateAvailabliliy,
					},
				}

				arrayLock.Lock()
				stations = append(stations, station)
				arrayLock.Unlock()
			}
			return nil
		})
	}

	if err := eGroup.Wait(); err != nil {
		return nil, err
	}

	return stations, nil
}
