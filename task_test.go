package goticks

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/parametalol/curry/assert"
	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
)

func TestTask(t *testing.T) {
	t.Run("collect ticks", func(t *testing.T) {
		ticker := ticker.New[int]()

		var ticks []int
		NewTask(ticker, func(tick int) {
			ticks = append(ticks, tick)
		}).Start()

		ticker.Tick(1).Wait()
		ticker.Tick(10).Wait()
		ticker.Tick(101).Wait()

		assert.That(t,
			assert.EqualSlices([]int{1, 10, 101}, ticks))
	})

	t.Run("task stop and start", func(t *testing.T) {
		ticker := ticker.New[int]()

		var ticks []int
		task := NewTask(ticker, func(tick int) {
			ticks = append(ticks, tick)
		})
		task.Start()
		ticker.Tick(1).Wait()
		task.Stop()
		ticker.Tick(10).Wait()
		task.Start()
		ticker.Tick(101).Wait()
		assert.That(t,
			assert.EqualSlices([]int{1, 101}, ticks))
	})

	t.Run("ont ticker, three tasks", func(t *testing.T) {
		ticker := ticker.New[int32]()

		var i atomic.Int32
		for range 3 {
			NewTask(ticker, func(tick int32) {
				i.Add(tick)
			}).Start()
		}
		ticker.Tick(1).Wait()
		ticker.Tick(10).Wait()
		ticker.Tick(101).Wait()

		assert.That(t,
			assert.Equal(int32(3*(1+10+101)), i.Load()))
	})
}

func Test_options(t *testing.T) {
	t.Run("on start", func(t *testing.T) {
		ticker := ticker.New[int]()

		var ticks []int
		task := NewTask(ticker, func(tick int) {
			ticks = append(ticks, tick)
		}, WithOnStart(func() error {
			ticks = append(ticks, 0)
			return errors.New("that's ok")
		}), WithOnStop(func() {
			ticks = append(ticks, -1)
		}),
		)

		task.Start()
		task.Start()

		ticker.Tick(1).Wait()
		ticker.Tick(10).Wait()
		ticker.Tick(101).Wait()

		task.Stop()

		assert.That(t,
			assert.EqualSlices([]int{0, 1, 10, 101, -1}, ticks))
	})

	t.Run("error on start", func(t *testing.T) {
		ticker := ticker.New[int]()

		var ticks []int
		task := NewTask(ticker, func(tick int) {
			ticks = append(ticks, tick)
		}, WithOnStart(func() error {
			ticks = append(ticks, 0)
			return utils.ErrStopped
		}), WithOnStop(func() {
			ticks = append(ticks, -1)
		}),
		)

		task.Start()

		ticker.Tick(1).Wait()
		ticker.Tick(10).Wait()
		ticker.Tick(101).Wait()

		task.Stop()

		assert.That(t,
			assert.EqualSlices([]int{0}, ticks))
	})

	t.Run("WithTickerStop", func(t *testing.T) {
		ticker := ticker.NewTimer(time.Second)

		task := NewTask(ticker, func() {},
			WithTickerStop())

		var ticks []time.Time
		tickerTicks := ticker.Ticks()
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			ticks = slices.Collect(tickerTicks)
			wg.Done()
		}()
		task.Start() // one tick on start.
		time.Sleep(500 * time.Millisecond)
		task.Stop() // stop before the first timed tick.

		time.Sleep(2 * time.Second)

		wg.Wait()
		assert.That(t,
			assert.Equal(1, len(ticks)))
	})

	t.Run("task stop and start WithTickerStop", func(t *testing.T) {
		ticker := ticker.New[int]()

		var ticks []int
		task := NewTask(ticker, func(tick int) {
			ticks = append(ticks, tick)
		}, WithTickerStop())
		task.Start()
		ticker.Tick(1).Wait()
		task.Stop()
		ticker.Tick(10).Wait()
		task.Start()
		ticker.Tick(101).Wait()
		assert.That(t,
			assert.EqualSlices([]int{1, 101}, ticks))
	})
}
