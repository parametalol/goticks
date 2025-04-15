package loop

import (
	"context"
	"errors"
	"iter"

	"github.com/parametalol/go-ticks/utils"
)

// OnTick calls task on every tick from the ticker.
// The function returns the last task error when the ticker is stopped, or task
// fails with [ErrStopped].
func OnTick[TickType any](ticks iter.Seq[TickType], task func(context.Context, TickType) error) error {
	ctx, cancel := context.WithCancelCause(context.Background())
	defer cancel(utils.ErrStopped)
	var err error
	for tick := range ticks {
		if err = task(ctx, tick); errors.Is(err, utils.ErrStopped) {
			// This returns false to the ticks iterator.
			break
		}
	}
	return err
}
