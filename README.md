# Go-ticks

Go-ticks is a lightweight library for managing periodic tasks in your Go applications.

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

What could be improved:

- Separation of concerns:
  - a ticker: the tick generator;
  - a task runner: the middleware that executes the task on ticker ticks;
  - a task: the code to be executed by the runner.
- Graceful termination: the task has to be provided with a context, that is cancelled by the controller demand.
- Task failure handling: the task should be able to notify the task runner about a permanent error.
- The ticker has to be stoppable and restartable.
- The first tick should arrive on start (not after first period).

This library assembles a set of types, that implement the above.

## Features

- Simple and intuitive interface for scheduling periodic tasks:
  - `Start()` for (re-)starting the periodic execution;
  - `Stop()` to gracefully interrupt the execution by cancelling the context;
  - `Wait()` to wait for the running tasks to terminate;
  - `Error()` to consult the termination reason.
- A `Ticker` interface and an implementation that ticks on start.
- A `TestTicker` implementation of the `Ticker` interface that allows for sending a code controlled ticks.
- A list of task wrappers, such as `NoOverlap`, `WithRetry` and others.

## Example

The example shows a periodic task, that prints "tick" 3 times once a second.

```go
package main

import (
  "fmt"
  "time"

  "github.com/parametalol/periodic"
)

func tick() { fmt.Println("tick") }

func main() {
  task := periodic.NewTask("tick task", time.Second, tick)
  task.Start()
  time.Sleep(2500 * time.Millisecond)
  task.Stop()
  task.Wait()
}
```
