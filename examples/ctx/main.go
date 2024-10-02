package main

import (
	"context"
	"fmt"
	"time"

	"github.com/dataphos/lib-shutdown/pkg/graceful"
)

func main() {
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
}
