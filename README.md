# learning-go

A collection of Go projects written by someone who definitely knew what they were doing from the start and did not rewrite anything three times.

## Projects

### 1. GitHub Activity CLI
Fetches GitHub events and prints them. Learned that Go has interfaces. Used one. Very proud.

### 2. Site Health Checker
Checks if websites are up. Started with a mutex. Switched to channels. Switched to a worker pool. README was rewritten to match the code that actually exists.

### 3. Flash Sale
100,000 goroutines fight over 100 items. Discovered that:
- unbuffered channels deadlock
- closing channels you're still sending to panics
- atomic load + atomic add is not atomic
- `sync.Once` exists and is great
- none of this was obvious until it exploded

### 4. Fan-In Fan-Out
1,000,000 products enriched in parallel, fanned out to 10 validator workers, fanned back in to one merged channel, cancelled by context after ~50 products because 10 workers × 1s sleep × 1M products would take 27 hours. Learned that:
- `[]chan T` and `[]<-chan T` are not interchangeable even though `chan T` and `<-chan T` are
- goroutines ranging over a shared channel is fan-out, not a bug
- `defer close(ch)` per worker is cleaner than a separate closer goroutine
- fan-in is just the reverse: one goroutine per input channel, all feeding one output
- `select` outside a loop checks context once at start and then ignores it forever, which is not cancellation
- `wg.Done()` must be deferred, not placed only in the happy path

## Progress

Learning Go by writing things, running them, watching them panic, and figuring out why. Working as intended.
