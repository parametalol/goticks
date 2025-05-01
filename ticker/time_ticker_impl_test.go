package ticker

import (
	"slices"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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

func TestNewTimer(t *testing.T) {
	timer := NewTimer(time.Second)
	assert.False(t, timer.(*timeTickerImpl).running.Load())

	time.AfterFunc(2500*time.Millisecond, timer.Stop)

	ticks := timer.Ticks()
	assert.True(t, timer.(*timeTickerImpl).running.Load())

	times := slices.Collect(ticks)
	assert.False(t, timer.(*timeTickerImpl).running.Load())

	if len(times) != 3 {
		t.Errorf("i expected to be %d, got %d", 3, len(times))
	}
}
