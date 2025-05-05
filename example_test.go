package goticks

import (
	"context"
	"fmt"
	"os"
	"time"

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
		utils.Retry[int](utils.SimpleRetryPolicy(3),
			utils.Log[int](os.Stdout, os.Stdout, "example",
				counter))).
		// Start the ticker loop in a goroutine:
		Start()

	for tick := range 10 {
		ticker.Tick(tick).
			// Wait for the tick to be processed
			// to ensure stable sequential output:
			Wait()
	}

	// Output:
	// Calling example
	// tick # 0
	// Calling example
	// tick # 1
	// Calling example
	// tick # 2
	// Execution of example failed after the first attempt with error: non-stop error
	// Retry 1 of example
	// tick # 2
	// Execution of example failed after retry 1 with error: non-stop error
	// Retry 2 of example
	// tick # 2
	// Execution of example failed after retry 2 with error: non-stop error
	// Calling example
	// tick # 3
	// Execution of example stopped after the first attempt with error: stop error: stopped
}

func ExampleTask_Stop() {
	ticker := ticker.New[int]()
	task := NewTask(ticker,
		func(tick int) {
			fmt.Println("Tick:", tick)
		})

	i := 0
	sendTicks := func() {
		for range 3 {
			ticker.Tick(i).Wait()
			i++
		}
	}

	task.Start()
	sendTicks()
	task.Stop()

	// These ticks are ignored by the task:
	sendTicks()

	task.Start()
	sendTicks()
	task.Stop()

	// Output:
	// Tick: 0
	// Tick: 1
	// Tick: 2
	// Tick: 6
	// Tick: 7
	// Tick: 8
}

func ExampleTask_Start() {
	startTime := time.Now()
	ticker := ticker.NewTimer(time.Second)
	NewTask(ticker,
		func(t time.Time) {
			fmt.Println("Current time:", t.Sub(startTime).Round(time.Second))
		}).Start()

	time.Sleep(3*time.Second + 10*time.Millisecond)
	ticker.Stop()

	// Output:
	// Current time: 0s
	// Current time: 1s
	// Current time: 2s
	// Current time: 3s
}
