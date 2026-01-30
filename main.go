package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mikesullivan63/downloader/downloader"
)

func main() {
	var out string
	flag.StringVar(&out, "o", "", "Output file (default: stdout)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "usage: downloader [-o output] <url>")
		os.Exit(2)
	}
	url := flag.Arg(0)
	var w *os.File
	var err error
	if out == "" {
		w = os.Stdout
	} else {
		w, err = os.Create(out)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create output file: %v\n", err)
			os.Exit(1)
		}
		defer w.Close()
	}
	if err := downloader.Download(url, w); err != nil {
		fmt.Fprintf(os.Stderr, "download failed: %v\n", err)
		os.Exit(1)
	}
}
