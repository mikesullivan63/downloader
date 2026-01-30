package downloader

import (
	"errors"
	"io"
	"net/http"
)

// Download fetches the contents at url and writes them to w.
func Download(url string, w io.Writer) error {
	if url == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("non-2xx response: " + resp.Status)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}
