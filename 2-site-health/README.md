# Site Health Checker

Checks HTTP status of URLs using a worker pool.

## What it covers

- Worker pool pattern: fixed `N` goroutines consume from a `jobs` channel
- Buffered channels for `jobs` and `results` to decouple producer/consumer
- `sync.WaitGroup` to detect when all workers finish, then close `results`
- Range over closed channel to drain results

## Usage

Add URLs (one per line) to `websites.txt`, then:

```
go run .
```

## How it works

```
main → fills jobs channel → workers consume jobs → push to results
                                       ↓
                             wg.Wait() → close(results) → drain + print
```

`workerCount = 3` controls parallelism.
