# Flash Sale

Simulates a flash sale: 100,000 users race to buy 100 items, with a 20% payment failure rate.

## What it covers

- Buffered channel as a resource pool (`itemsInStock`)
- `done` channel for clean shutdown signalling across goroutines
- `sync.Once` to safely close `done` exactly once
- `atomic.Int32` for lock-free sale counter
- `select` for racing between "claim item" and "sale over"
- Returning items to pool on payment failure without panic

## How it works

```
main → fills itemsInStock channel (100 slots)
     → spawns 100k goroutines

each goroutine:
  select {
  case item := <-itemsInStock:  claim slot
      paymentGateRoll() → 80% pass, 20% fail
      pass → itemsBought++; if 100th sale → close(done)
      fail → return item to channel; exit
  case <-done:                  sale over, exit
  }
```

`paymentGateRoll` sleeps 1s to simulate network latency.

## Usage

```
go run .
```
