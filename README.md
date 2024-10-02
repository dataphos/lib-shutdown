# lib-shutdown

lib-shutdown is a lightweight Go library that offers utility functions that deal with graceful termination on selected
OS signals.

## Installation

`go get github.com/dataphos/lib-shutdown`

## Getting Started

### Cancelling Context On OS Signals

The code snippet below shows the basic principle of the `graceful` package:

```go
secondsRemaining := 10
ctx, cancel := context.WithTimeout(context.Background(), time.Duration(secondsRemaining)*time.Second)
defer cancel()

ctx = graceful.WithSignalShutdown(ctx)

for {
    select {
    case <-ctx.Done():
        fmt.Println("context cancelled, leaving")
        return
    case <-time.After(1 * time.Second):
        secondsRemaining -= 1
        fmt.Println(secondsRemaining, "seconds remaining")
    }
}
```

`graceful.WithSignalShutdown` returns a new ctx derived from the one given, which will be cancelled when this process
receives a SIGTERM or SIGQUIT signal, cancelling the returned context before the parent context times out.

### Graceful Shutdown of HTTP Servers

`graceful.ListenAndServe` and `graceful.ListenAndServeTLS` functions wrap the equivalent
method of the given `http.Server` instance, signaling to the server that a SIGTERM or SIGQUIT signal
was received, allowing some additional time for the server to cleanup and complete the active requests.
