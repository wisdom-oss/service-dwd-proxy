package internal

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gabriel-vasile/mimetype"
)

const filenamePattern = "dwd-proxy-*%s"

var errStatusNot200 = errors.New("the remote server did not indicate a successful request")

func Download(uri string) (filepath string, err error) {
	res, err := http.Get(uri) //nolint:gosec
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", errStatusNot200
	}

	mime := mimetype.Lookup(res.Header.Get("Content-Type"))
	if mime == nil {
		mime = &mimetype.MIME{}
	}

	filename := fmt.Sprintf(filenamePattern, mime.Extension())

	f, err := os.CreateTemp("", filename)
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(f, res.Body); err != nil {
		return "", err
	}

	if err := f.Sync(); err != nil {
		return "", err
	}

	if err := f.Close(); err != nil {
		return "", err
	}

	return f.Name(), nil
}
