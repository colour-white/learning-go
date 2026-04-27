# Fan-In Fan-Out

Simulates a product enrichment and validation pipeline: 1,000,000 products enriched in parallel, fanned out to 10 validator workers, fanned back in to a single merged channel. With context cancellation, because apparently just letting goroutines run forever is frowned upon.

## What it covers

- Fan-out: single input channel split across multiple worker goroutines
- Fan-in: multiple output channels merged into one
- `sync.WaitGroup` to close channels only after all workers finish
- Buffered channels to avoid blocking producers
- Receive-only (`<-chan`) and send-only (`chan<-`) channel types
- `defer close(ch)` per worker for clean shutdown
- `context.Context` propagated through the entire pipeline, used correctly on the third attempt

## How it works

```
EnrichedProducts:
  1,000,000 goroutines → each enriches one Product via faker
  select: ctx.Done() → exit, or send to buffered channel
  closer goroutine waits on WaitGroup → closes channel

TrustValidate (fan-out):
  10 worker goroutines compete for products from input channel
  per product: select inside loop - ctx.Done() → return, or validate and send
  defer close per worker when input exhausted

MergeChannels (fan-in):
  10 goroutines → each drains one worker channel
  per item: select inside loop - ctx.Done() → return, or forward to result channel
  closer goroutine waits on WaitGroup → closes merged channel

main:
  5s timeout context
  EnrichedProducts → TrustValidate → MergeChannels → for range merged → print
```

## Caveats

`EnrichedProducts` spawns one goroutine per product. At 1,000,000 products this launches 1M goroutines simultaneously, which the Go runtime handles with surprising grace and the machine handles with surprising heat. The fix - a worker pool, like `TrustValidate` already uses - was left as an exercise in consistency.

The pipeline processes roughly 10 products/second (10 workers × 1s sleep per product). The timeout is 25 seconds. At 1M products this means context cancels after ~250 items, the pipeline drains cleanly, and the program exits having enriched 1M products and validated roughly 0.025% of them. Very thorough.

Context cancellation was added in stages: first missing `wg.Done()` on cancel, then select outside the loop, then select inside the wrong place. The current version is correct. Probably.

## Usage

```
go run .
```
