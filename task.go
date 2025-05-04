package goticks

import (
	"context"
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

	once   sync.Once
	paused atomic.Bool
}

var _ Task = (*taskImpl[any])(nil)

func NewTask[TickType any, Fn utils.Func[TickType]](ticker ticker.Tickable[TickType], fn Fn) ticker.Restartable {
	task := &taskImpl[TickType]{
		ticker: ticker,
	}
	adaptedTask := utils.Adapt[TickType](fn)
	task.task = func(ctx context.Context, tick TickType) error {
		if task.paused.Load() {
			return nil
		}
		return adaptedTask(ctx, tick)
	}
	return task
}

// Start another task execution loop.
func (t *taskImpl[TickType]) Start() {
	t.paused.Store(false)
	t.once.Do(func() {
		go func() {
			_ = loop.OnTick(t.ticker.Ticks(), t.task)
		}()
	})
}

// Stop all running loops by stopping the ticker.
func (t *taskImpl[TickType]) Stop() {
	t.paused.Store(true)
}
