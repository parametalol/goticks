package ticker

import (
	"fmt"
	"iter"
	"time"
)

func Example_tickerImpl_Tick() {
	ticker := New[int]()
	defer ticker.Stop()

	ticks := ticker.Ticks()

	ticker.Tick(42)
	next, stop := iter.Pull(ticks)

	fmt.Println(next())
	stop()

	// Output:
	// 42 true
}

// This example illustrates the use of the ticker, that ticks every second
// during 2.5 seconds, and stops.
func ExampleNewTimer() {

	timer := NewTimer(time.Second)
	time.AfterFunc(2500*time.Millisecond, timer.Stop)

	startTime := time.Now()
	for tick := range timer.Ticks() {
		fmt.Println(tick.Sub(startTime).Round(time.Second))
	}

	// Output:
	// 0s
	// 1s
	// 2s
}
