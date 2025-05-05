package goticks

import (
	"context"
	"errors"
	"sync"
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

	once    sync.Once
	started atomic.Bool
}

var _ Task = (*taskImpl[any])(nil)

func NewTask[TickType any, Fn utils.Func[TickType]](ticker ticker.Tickable[TickType], fn Fn, opts ...option) ticker.Restartable {
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

// Start another task execution loop.
func (t *taskImpl[TickType]) Start() {
	if t.started.Load() {
		return
	}
	if t.options.onStart != nil && errors.Is(t.options.onStart(), utils.ErrStopped) {
		return
	}
	t.started.Store(true)
	t.once.Do(func() {
		ticks := t.ticker.Ticks()
		go func() {
			_ = loop.OnTick(ticks, t.task)
		}()
	})
}

// Stop all running loops by stopping the ticker.
func (t *taskImpl[TickType]) Stop() {
	if t.started.Swap(false) && t.options.onStop != nil {
		go t.options.onStop()
	}
}
