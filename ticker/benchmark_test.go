package ticker

import (
	"testing"
)

// BenchmarkTicker_TickWait measures the overhead of sending and acknowledging ticks.
func BenchmarkTicker_TickWait(b *testing.B) {
	t := New[int]()
	// Start a consumer that drains ticks without work.
	seq := t.Ticks()
	go func() {
		for range seq {
		}
	}()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t.Tick(i).Wait()
	}
}
