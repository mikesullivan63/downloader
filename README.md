# downloader

Simple CLI to download a URL to a file using the standard library.

Usage

- Build: `go build ./...`
- Run: `go run . -o out.html https://example.com`
- Test: `go test ./...`

Files

- `main.go` - CLI wrapper that calls `downloader.Download`
- `downloader/downloader.go` - download logic and tests
