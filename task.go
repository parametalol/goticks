package goticks

import (
	"context"

	"github.com/parametalol/goticks/loop"
	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
)

type Task interface {
	Start()
	Stop()
}

type taskImpl[TickType any] struct {
	ticker ticker.Ticker[TickType]
	task   func(context.Context, TickType) error
}

var _ Task = (*taskImpl[any])(nil)

func NewTask[TickType any, Fn utils.Func[TickType]](ticker ticker.Ticker[TickType], task Fn) ticker.Restartable {
	return &taskImpl[TickType]{
		ticker: ticker,
		task:   utils.Adapt[TickType](task),
	}
}

// Start another task execution loop.
func (t *taskImpl[TickType]) Start() {
	go loop.OnTick(t.ticker.Ticks(), t.task)
}

// Stop all running loops by stopping the ticker.
func (t *taskImpl[TickType]) Stop() {
	t.ticker.Stop()
}
