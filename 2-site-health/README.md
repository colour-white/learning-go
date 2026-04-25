# Site Health Checker

A Go CLI tool that concurrently checks if a list of websites are up or down.

## What it covers

- Concurrent HTTP requests with goroutines
- `sync.WaitGroup` — waiting for all goroutines to finish
- `sync.Mutex` — safe concurrent map writes
- HTTP client with timeout (`http.Client{Timeout: ...}`)
- File reading + string splitting

## Usage

Add URLs (one per line) to `websites.txt`, then:

```
go run .
```

Output:
```
https://google.com - ok
https://somedownsite.com - down
```
