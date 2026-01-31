package main

import (
	"flag"
	"fmt"
	"os"
	"time"
	"math/rand"
	"sync"
	"runtime"

	"github.com/mikesullivan63/downloader/downloader"
	"github.com/mikesullivan63/downloader/messages"
)

// Type aliases for message structs
type PageDiscovered = messages.PageDiscovered
type ImageDiscovered = messages.ImageDiscovered
type JobStatus = messages.JobStatus

// New worker-pool implementation (keeps the old malformed functions untouched)
func startWorkerPool2(enqueue chan PageDiscovered, pageChan chan PageDiscovered, imageChan chan ImageDiscovered, status *JobStatus, wg *sync.WaitGroup, workers int) chan struct{} {
	forwarderDone := make(chan struct{})
	go func() {
		for p := range enqueue {
			pageChan <- p
		}
		close(pageChan)
		close(forwarderDone)
	}()

	for i := 0; i < workers; i++ {
		go func() {
			for ev := range pageChan {
				if err := downloader.Scan(ev, enqueue, imageChan, status, wg); err != nil {
					fmt.Fprintf(os.Stderr, "scanner error: %v\n", err)
				}

				fmt.Printf("Finished page: %+v\n%+v\n", ev, status)

				time.Sleep(1 * time.Second)

				wg.Done()
			}
		}()
	}

	return forwarderDone
}

func startImagePrinter2(imageChan chan ImageDiscovered, done chan struct{}, status *JobStatus) {
	go func() {
		for img := range imageChan {
//			fmt.Printf("found image: %s\n", img.URL)
			if err := downloader.DownloadImage(img.URL, status); err != nil {
				fmt.Fprintf(os.Stderr, "download image error: %v\n", err)
			}

			fmt.Printf("Finished image: %+v\n%+v\n", img, status)
		}
		close(done)
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
	enqueueChan := make(chan PageDiscovered, 100)
	pageChannel := make(chan PageDiscovered)
	imageChannel := make(chan ImageDiscovered, 100)

	// jobCompleted := make(chane int)

	var wg sync.WaitGroup
	workers := runtime.NumCPU()

	forwarderDone := startWorkerPool2(enqueueChan, pageChannel, imageChannel, &status, &wg, workers)

	// ensure the initial seed is counted before waiting
	wg.Add(1)

	imageDone := make(chan struct{})
	startImagePrinter2(imageChannel, imageDone, &status)

	// seed with initial page (enqueue increments waitgroup via forwarder)
	enqueueChan <- PageDiscovered{JobID: jobID, Depth: 0, URL: url}

	// wait for all pages to be processed
	wg.Wait()

	fmt.Printf("All pages done: %+v\n", status)

	// no more enqueues, close to let forwarder finish
	close(enqueueChan)
	<-forwarderDone

	// close images and wait for the image printer to finish
	close(imageChannel)
	<-imageDone

	fmt.Printf("Done. status: %+v\n", status)
}

