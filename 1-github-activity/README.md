# GitHub Activity CLI

A Go CLI tool that fetches and displays GitHub user info and public event history.

## What it covers

- HTTP requests with custom headers (`net/http`)
- JSON decoding of API responses (`encoding/json`)
- Interface design - `Printable` interface implemented by multiple types
- CLI flags (`flag` package)
- Environment config via `.env` (`godotenv`)
- Error propagation - returning errors up the call stack vs `log.Fatal`

## Usage

```
go run main.go                  # text output
go run main.go -json            # JSON output
```

Requires a `.env` file with `GIT_TOKEN=<your_github_token>`.
