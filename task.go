package goticks

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/parametalol/goticks/loop"
	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
)

type Task interface {
	Start()
	Stop()
}

type taskImpl[TickType any] struct {
	ticker ticker.Tickable[TickType]
	task   func(context.Context, TickType) error

	options options

	once    atomic.Bool
	started atomic.Bool
}

var _ Task = (*taskImpl[any])(nil)

type RestartableWithTicker[TickType any] interface {
	ticker.Restartable
	Ticker() ticker.Tickable[TickType]
}

// NewTask returns an instance of a restartable task, executed on the ticker
// ticks.
//
// The execution of tasks is paused on [Stop] and resumed on [Start] without
// affecting the ticker unless [WithTickerStop] is provided.
//
// If [WithTickerStop] is provided as an option, the ticker will be stopped on
// [Stop], which will interrupt all current ticks consumers. It will also be
// started on [Start], but the previously stopped consumers, except the current
// task, will not restart.
//
// Example:
//
//	NewTask(ticker.NewTimer(time.Second), task).Start() // run task every second
func NewTask[TickType any, Fn utils.Func[TickType]](ticker ticker.Tickable[TickType], fn Fn, opts ...option) RestartableWithTicker[TickType] {
	task := &taskImpl[TickType]{
		ticker: ticker,
	}
	for _, opt := range opts {
		opt(&task.options)
	}
	adaptedTask := utils.Adapt[TickType](fn)
	task.task = func(ctx context.Context, tick TickType) error {
		if !task.started.Load() {
			return nil
		}
		return adaptedTask(ctx, tick)
	}
	return task
}

// Start the task execution loop, once.
func (t *taskImpl[TickType]) Start() {
	if t.started.Swap(true) {
		return
	}
	if t.options.onStart != nil && errors.Is(t.options.onStart(), utils.ErrStopped) {
		t.started.Store(false)
		return
	}
	if !t.once.Swap(true) {
		ticks := t.ticker.Ticks()
		go func() {
			_ = loop.OnTick(ticks, t.task)
		}()
	}
}

// Stop all running loops by stopping the ticker.
func (t *taskImpl[TickType]) Stop() {
	if !t.started.Swap(false) {
		return
	}
	if t.options.stopTicker {
		if ticker, isStoppable := t.ticker.(ticker.Stoppable); isStoppable {
			ticker.Stop()
			t.once.Store(false)
		}
	}
	if t.options.onStop != nil {
		t.options.onStop()
	}
}

// Ticker returns the ticker, used for the task initialization.
func (t *taskImpl[TickType]) Ticker() ticker.Tickable[TickType] {
	return t.ticker
}
