# Goticks

Goticks is a lightweight library for managing periodic tasks in your Go applications.

## Rationale

Consider this example from the standard library, that prints current time every second during 10 seconds period:

```go
func ExampleNewTicker() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
		}
	}
}
```

The following might be needed to build a real use case:

- Concerns should be better separated to:
  - a task: the business logic to be executed by the runner;
  - a ticker: the tick generator, replacable with a different implementation;
  - a task runner: the middleware that executes the task on ticker ticks, logs the execution
    process, retries on failure, etc.
- Graceful termination: the task has to be provided with a cancellable context.
- Task failure handling: the task should be able to notify the task runner about a permanent error.
- The ticker has to be stoppable and restartable.
- The first tick should arrive on start (not after first period).

This library assembles a set of types, that implement the above.

## Example

```go
func ExampleNewTask() {
	ticker := ticker.NewTimer(time.Second)
	startTime := time.Now()

	NewTask(ticker, func(t time.Time) {
		fmt.Println("Current time:", t.Sub(startTime).Round(time.Second))
	}).Start()

	time.Sleep(3 * time.Second)
	ticker.Wait()
	ticker.Stop()

	// Output:
	// Current time: 0s
	// Current time: 1s
	// Current time: 2s
	// Current time: 3s
}
```
