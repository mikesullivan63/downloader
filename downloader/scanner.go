package downloader

import (
	"errors"
	"net/http"
	"log"
	"github.com/mikesullivan63/downloader/messages"
	"github.com/PuerkitoBio/goquery"
)

// Type aliases for message structs
type PageDiscovered = messages.PageDiscovered
type ImageDiscovered = messages.ImageDiscovered

const MAX_DEPTH = 3

// Scans page and gets links and images.
func Scanner(event PageDiscovered, pageChannel chan PageDiscovered, imageChannel chan ImageDiscovered) error {

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


	_, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing %s: %v", event.URL, err)
		return err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Printf("Error parsing %s: %v", event.URL, err)
		return err
	}

	if event.Depth < MAX_DEPTH {
		doc.Find("a").Each(func(i int, s *goquery.Selection) { 
			text, exists := s.Attr("href")

			if(exists) {
				pageEvent := PageDiscovered{JobID: event.JobID, Depth: event.Depth + 1, URL: text}
				pageChannel <- pageEvent
			}
		})
	}

	doc.Find("img").Each(func(i int, s *goquery.Selection) { 
		text, exists := s.Attr("src")

		if(exists) {
			imageEvent := ImageDiscovered{JobID: event.JobID, URL: text}
			imageChannel <- imageEvent
		}
	})


	return nil
}
