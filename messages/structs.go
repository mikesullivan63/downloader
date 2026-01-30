package messages

import (
	"time"
)


// Message types used by the downloader
type PageDiscovered struct {
	JobID int
	Depth int
	URL   string
}

type ImageDiscovered struct {
	JobID int
	URL   string
}

type JobStatus struct {
	JobID int
	PagesFound int
	PagesScanned int
	ImagesFound int
	ImagesScanned int
	LastActivity time.Time
}

