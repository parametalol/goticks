package loop

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
	"github.com/stretchr/testify/assert"
)

type tickerWithTick[TickType any] interface {
	ticker.Ticker[TickType]
	Tick(TickType) ticker.Waitable
}

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

	t.Run("one ticker two loops", func(t *testing.T) {
		var arr []int
		collector := func(tick int) {
			arr = append(arr, tick)
		}
		ticker := ticker.New[int]()
		mux := &sync.Mutex{}
		for range 3 {
			go OnTick(ticker.Ticks(), utils.Sync[int](mux, collector))
		}
		for tick := range 3 {
			ticker.Tick(tick).Wait()
		}
		assert.Equal(t, []int{0, 0, 0, 1, 1, 1, 2, 2, 2}, arr)
	})
}

func tickInRange(ticker tickerWithTick[int], n int) {
	for tick := range n {
		ticker.Tick(tick)
	}
	ticker.Wait()
	ticker.Stop()
}
