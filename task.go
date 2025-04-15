package go_ticks

import (
	"context"

	"github.com/parametalol/go-ticks/loop"
	"github.com/parametalol/go-ticks/ticker"
	"github.com/parametalol/go-ticks/utils"
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

func NewTask[TickType any, Fn utils.Func[TickType]](ticker ticker.Ticker[TickType], task Fn) *taskImpl[TickType] {
	return &taskImpl[TickType]{
		ticker: ticker,
		task:   utils.Adapt[TickType](task),
	}
}

func (t *taskImpl[TickType]) Start() {
	go loop.OnTick(t.ticker.Ticks(), t.task)
}

func (t *taskImpl[TickType]) Stop() {
	if tt, ok := t.ticker.(ticker.TimeTicker); ok {
		tt.Reset(0)
	}
}
