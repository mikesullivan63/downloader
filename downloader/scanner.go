package downloader

import (
	"errors"
	"fmt"
	"net/http"
	"log"
	"strings"
	"sync"
	"time"
	"github.com/mikesullivan63/downloader/messages"
	"github.com/PuerkitoBio/goquery"
)

// Type aliases for message structs
type PageDiscovered = messages.PageDiscovered
type ImageDiscovered = messages.ImageDiscovered
type JobStatus = messages.JobStatus

const MAX_DEPTH = 1

func Cleanup(status *JobStatus)  {
	// done scanning this page
	status.PagesScanned++
	status.LastActivity = time.Now()
}
 
// Scans page and gets links and images.
func Scan(event PageDiscovered, enqueue chan<- PageDiscovered, imageChannel chan<- ImageDiscovered, status *JobStatus, wg *sync.WaitGroup) error {
	fmt.Printf("Scanning page: %+v\n%+v\n", event, status)

	// mark that we're scanning this page
	defer Cleanup(status)

	if event.URL == "" {
		return errors.New("url is empty")
	}

	resp, err := http.Get(event.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("non-2xx response: " + resp.Status)
	}


	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing %s: %v", event.URL, err)
		return err
	}

	if event.Depth < MAX_DEPTH {
		doc.Find("a").Each(func(i int, s *goquery.Selection) {
			text, exists := s.Attr("href")

			if exists {
				pageEvent := PageDiscovered{JobID: event.JobID, Depth: event.Depth + 1, URL: text}
				status.PagesFound++
				wg.Add(1)
				enqueue <- pageEvent
			}
		})
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		text, exists := s.Attr("src")

		if exists && strings.HasPrefix(text, "http") {

			imageEvent := ImageDiscovered{JobID: event.JobID, URL: text}
			status.ImagesFound++
			imageChannel <- imageEvent
		}
	})

	return nil
}
