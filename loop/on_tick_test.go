package loop

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/parametalol/go-ticks/ticker"
	"github.com/parametalol/go-ticks/utils"
	"github.com/stretchr/testify/assert"
)

func TestOnTick(t *testing.T) {
	t.Run("loop on 3 test ticks", func(t *testing.T) {
		ticker := ticker.New[int]()
		ticks := ticker.Ticks()
		var i atomic.Int32

		go tickInRange(ticker, 3)

		err := OnTick(ticks, utils.Adapt[int](func(tick int) {
			i.Add(int32(tick))
		}))

		assert.Nil(t, err)
		assert.Equal(t, int32(3), i.Load())
	})

	t.Run("failing function", func(t *testing.T) {
		ticker := ticker.New[int]()
		ticks := ticker.Ticks()

		go tickInRange(ticker, 3)

		errTest := errors.New("test")
		err := OnTick(ticks, func(context.Context, int) error {
			return errTest
		})
		assert.ErrorIs(t, err, errTest)
	})

	t.Run("error handling", func(t *testing.T) {
		tempErr := errors.New("non-stop error")
		permErr := fmt.Errorf("stop error: %w", utils.ErrStopped)

		counter := func(_ context.Context, tick int) error {
			switch tick {
			case 2:
				return tempErr
			case 3:
				return permErr
			}
			return nil
		}

		ticker := ticker.New[int]()
		ticks := ticker.Ticks()

		go tickInRange(ticker, 5)

		err := OnTick(ticks, counter)
		assert.ErrorIs(t, err, utils.ErrStopped)

	})
}

func tickInRange(ticker ticker.TickerWithTick[int], n int) {
	for tick := range n {
		ticker.Tick(tick)
	}
	ticker.Wait()
	ticker.Stop()
}
