package downloader

import (
	"errors"
	"io"
	"fmt"
	"math/rand"
	"net/http"
	//"net/url"
	"os"
	"time"
)

func CleanupImageCount(status *JobStatus)  {
	// done scanning this page
	status.ImagesScanned++
	status.LastActivity = time.Now()
}

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

// DownloadImage fetches an image and writes it to disk, updating status.
func DownloadImage(imageUrl string, status *JobStatus) error {
	defer CleanupImageCount(status)

	if imageUrl == "" {
		return errors.New("url is empty")
	}
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("non-2xx response: " + resp.Status)
	}

	var w *os.File

	imgPath := fmt.Sprintf("IMAGES/%d", rand.Intn(10000000))
	w, err = os.Create(imgPath)
	if err != nil {
		return err
	}
	defer w.Close()
	_, err = io.Copy(w, resp.Body)
	return err
}
