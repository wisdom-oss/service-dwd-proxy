package helpers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile downloads the file supplied in the url parameter and writes it
// into a temporary file.
// The file handle of the temporary file will be returned
func DownloadFile(url string) (*os.File, error) {
	// create a temporary file
	f, err := os.CreateTemp("", "dwd-weather-data-service-*")
	if err != nil {
		return nil, err
	}
	// use the http client to request the file
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// now check if status 200 was returned
	if res.StatusCode != 200 {
		return nil, errors.New("unable to download file if status code is non-200")
	}
	// since all checks passed, write the contents into the file
	_, err = io.Copy(f, res.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(f.Name())
	return f, nil
}
