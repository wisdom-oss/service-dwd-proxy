package dwd

import (
	"io"
	"net/http"
	"os"
)

func Download(url string) (string, error) {
	_ = os.MkdirAll("/tmp", os.ModeDir|os.ModePerm)
	f, err := os.CreateTemp("", "dwd-proxy-*")
	if err != nil {
		return "", err
	}
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if res.StatusCode != http.StatusOK {
		return "", ErrResponseNotOK
	}

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}

	err = f.Sync()
	if err != nil {
		return "", err
	}
	err = f.Close()
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}
