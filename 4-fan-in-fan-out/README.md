# Fan-In Fan-Out

Simulates a product enrichment and validation pipeline: 1000 products enriched in parallel, fanned out to 10 validator workers, fanned back in to a single merged channel.

## What it covers

- Fan-out: single input channel split across multiple worker goroutines
- Fan-in: multiple output channels merged into one
- `sync.WaitGroup` to close channels only after all workers finish
- Buffered channels to avoid blocking producers
- Receive-only (`<-chan`) and send-only (`chan<-`) channel types
- `defer close(ch)` per worker for clean shutdown

## How it works

```
EnrichedProducts:
  1000 goroutines → each enriches one Product via faker → sends to buffered channel
  closer goroutine waits on WaitGroup → closes channel

TrustValidate (fan-out):
  10 worker goroutines compete for products from input channel
  each worker → own output channel → validates, scores, promotes
  defer close per worker when input exhausted

MergeChannels (fan-in):
  10 goroutines → each drains one worker channel → forwards to single result channel
  closer goroutine waits on WaitGroup → closes merged channel

main:
  EnrichedProducts → TrustValidate → MergeChannels → print results
```

## Caveats

`EnrichedProducts` spawns one goroutine per product. Very elegant. Very scalable. Works great at 1000. At 1M it spawns 1M goroutines simultaneously, at which point the OS kindly introduces itself and the program dies. The fix - a worker pool like `TrustValidate` uses - was too boring to apply consistently.

## Usage

```
go run .
```
