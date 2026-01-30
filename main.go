package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"math/rand"

	"github.com/mikesullivan63/downloader/downloader"
	"github.com/mikesullivan63/downloader/messages"
)

// Type aliases for message structs
type PageDiscovered = messages.PageDiscovered
type ImageDiscovered = messages.ImageDiscovered
type JobStatus = messages.JobStatus

func startPageWorker(pageChannel chan PageDiscovered, imageChannel chan ImageDiscovered, status *JobStatus) {
	// consume events and start a scanner for each
	go func() {
		for ev := range pageChannel {
			// run scanner in its own goroutine so it can publish back to pageChannel
			ev := ev
			go func() {
				if err := downloader.Scan(ev, pageChannel, imageChannel, status); err != nil {
					fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
				}
			}()

			fmt.Printf("Finished page\n%+v\n", status)

			time.Sleep(1 * time.Second)
		}
	}()
}

func startImageWorker(imageChannel chan ImageDiscovered, status *JobStatus) {
	// consume events and start a scanner for each
	go func() {
		for ev := range imageChannel {
			// run scanner in its own goroutine so it can publish back to pageChannel
			ev := ev
			go func() {
				if err := downloader.Download(ev.URL, status); err != nil {
					fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
				}
			}()

			fmt.Printf("Finished image\n%+v\n", status)

		}
	}()
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: downloader <url>")
		os.Exit(2)
	}
	url := flag.Arg(0)

	jobID := rand.Intn(100)

	status := JobStatus{ JobID:jobID, PagesFound: 1, PagesScanned: 0, ImagesFound: 0, ImagesScanned: 0, LastActivity: time.Now() }
	pageChannel := make(chan PageDiscovered)
	imageChannel := make(chan ImageDiscovered)

	// start workers
	startPageWorker(pageChannel, imageChannel, &status)
	// startImageWorker(imageChannel, &status)

	/*
	// print found images
	go func() {
		for img := range imageChannel {
			fmt.Printf("found image: %s\n", img.URL)
		}
	}()
		*/

	// seed with initial page
	pageChannel <- PageDiscovered{JobID: jobID, Depth: 0, URL: url}

	// block forever for now (or until the process is killed)
	select {}
}

