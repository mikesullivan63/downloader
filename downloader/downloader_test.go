package downloader

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownload_Success(t *testing.T) {
	body := "hello world"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, body)
	}))
	defer ts.Close()
	var buf bytes.Buffer
	if err := Download(ts.URL, &buf); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != body {
		t.Fatalf("expected %q got %q", body, buf.String())
	}
}

func TestDownload_BadStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		io.WriteString(w, "oops")
	}))
	defer ts.Close()
	var buf bytes.Buffer
	err := Download(ts.URL, &buf)
	if err == nil {
		t.Fatalf("expected error for non-2xx response")
	}
}
