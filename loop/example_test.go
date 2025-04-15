package loop

import (
	"context"
	"fmt"
	"time"

	"github.com/parametalol/go-ticks/ticker"
)

// This example runs a loop, controlled by a periodic ticker.
// The called function is wrapped with a logger, that logs on every invocation
// and on error.
func ExampleOnTick() {
	ticker := ticker.NewTimer(time.Second)
	time.AfterFunc(2500*time.Millisecond, ticker.Stop)
	startTime := time.Now()

	// The tick function returns temporary error.
	err := OnTick(ticker.Ticks(),
		func(_ context.Context, tick time.Time) error {
			fmt.Println("tick", tick.Sub(startTime).Round(time.Second))
			return fmt.Errorf("oops")
		})

	fmt.Println(err)

	// Output:
	// tick 0s
	// tick 1s
	// tick 2s
	// oops
}
