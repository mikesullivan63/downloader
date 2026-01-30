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

// Download fetches the contents at url and writes them to w.
func Download(imageUrl string, status *JobStatus) error {
	
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

	//trimmed, err := url.Parse(imageUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse umage url: %v\n", err)
		os.Exit(1)
	}

//	w, err = os.Create("IMAGES/" + trimmed.Path)
	imgPath := fmt.Sprintf("IMAGES/%d",rand.Intn(10000000) )
	w, err = os.Create(imgPath )
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create output file: %v\n", err)
		os.Exit(1)
	}

	defer w.Close()

	_, err = io.Copy(w, resp.Body)

	// done scanning this page
	status.ImagesScanned++
	status.LastActivity = time.Now()

	return err
}
