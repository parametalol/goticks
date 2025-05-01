package ticker

import (
	"sync/atomic"
	"testing"
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
}
