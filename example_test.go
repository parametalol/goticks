package goticks

import (
	"context"
	"fmt"
	"os"

	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
)

func ExampleNewTask() {
	counter := func(ctx context.Context, tick int) error {
		fmt.Println("tick #", tick)
		switch tick {
		case 2:
			return fmt.Errorf("non-stop error")
		case 3:
			return fmt.Errorf("stop error: %w", utils.ErrStopped)
		}
		return nil
	}

	// This ticker issues ticks of type int.
	ticker := ticker.New[int]()

	// NewTask initializes a procedure to call counter on the ticker ticks, log
	// the attempts and the errors, and retry 2 times on non-permanent errors.
	NewTask(ticker,
		utils.WithRetry[int](utils.SimpleRetryPolicy(3),
			utils.WithLog[int](os.Stdout, os.Stdout, "example",
				counter))).
		// Start the process in a goroutine:
		Start()

	for tick := range 10 {
		ticker.Tick(tick).
			// Wait for the tick to be processed
			// to ensure stable sequential output:
			Wait()
	}

	// Let's wait for all currently running ticker senders to complete before
	// stopping the ticker.
	ticker.Wait()

	// Stop will cancel the context for the running tasks.
	ticker.Stop()

	// Output:
	// Calling example
	// tick # 0
	// Calling example
	// tick # 1
	// Calling example
	// tick # 2
	// Execution of example failed with error: non-stop error
	// Retry 1 of example
	// tick # 2
	// Execution of example failed after retry 1 with error: non-stop error
	// Retry 2 of example
	// tick # 2
	// Execution of example failed after retry 2 with error: non-stop error
	// Calling example
	// tick # 3
	// Execution of example stopped with error: stop error: stopped
}
