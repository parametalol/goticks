package ticker

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("full manual", func(t *testing.T) {
		ticker := New[int32]()
		ticks := ticker.Ticks()

		var i atomic.Int32
		go func() {
			for tick := range ticks {
				i.Add(tick)
			}
		}()

		for i := range 3 {
			ticker.Tick(int32(i)).Wait()
		}
		ticker.Stop()

		if i.Load() != 3 {
			t.Errorf("i expected to be %d, got %d", 3, i.Load())
		}
	})
	t.Run("full timer", func(t *testing.T) {
		timer := NewTimer(time.Second)
		time.AfterFunc(2500*time.Millisecond, timer.Stop)

		ticks := slices.Collect(timer.Ticks())

		if len(ticks) != 3 {
			t.Errorf("i expected to be %d, got %d", 3, len(ticks))
		}
	})
}

func TestTicker_Reset(t *testing.T) {
	ticker := NewTimer(0)
	ticker.Reset(time.Second)
	time.AfterFunc(2500*time.Millisecond, ticker.Stop)
	var i atomic.Int32
	for range ticker.Ticks() {
		i.Add(1)
	}

	if i.Load() != 3 {
		t.Errorf("i expected to be %d, got %d", 3, i.Load())
	}
}
