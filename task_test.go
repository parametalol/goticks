package goticks

import (
	"errors"
	"sync/atomic"
	"testing"

	"github.com/parametalol/goticks/ticker"
	"github.com/parametalol/goticks/utils"
	"github.com/stretchr/testify/assert"
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

		assert.Equal(t, []int{1, 10, 101}, ticks)
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
		assert.Equal(t, []int{1, 101}, ticks)
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

		assert.Equal(t, int32(3*(1+10+101)), i.Load())
	})

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

		assert.Equal(t, []int{0, 1, 10, 101, -1}, ticks)
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

		assert.Equal(t, []int{0}, ticks)
	})

}
