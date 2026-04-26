# Site Health Checker (Channel variant)

Same as `2-site-health` but uses a channel instead of a mutex-protected map to collect results.

## What it covers

- Goroutines + channels for producer/consumer pattern
- Closing channels to signal completion
- `sync.WaitGroup` in a separate goroutine to close channel while receiver is active
- Why unbuffered channels require concurrent sender and receiver

## Usage

Add URLs (one per line) to `websites.txt`, then:

```
go run .
```

## vs 2-site-health

| | 2-site-health | 2.1-site-health |
|---|---|---|
| Sync primitive | `sync.Mutex` + map | channel |
| Results | printed after all done | printed as they arrive |
