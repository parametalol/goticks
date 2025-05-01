package ticker

import (
	"iter"
	"sync"
	"sync/atomic"
	"time"
)

type timeTickerImpl struct {
	tickerImpl[time.Time]
	resetCh  chan time.Duration
	duration atomic.Int64

	running atomic.Bool
	runWg   sync.WaitGroup
}

var _ TimeTicker = (*timeTickerImpl)(nil)

// NewTimer creates a ticker that ticks on a timer.
// The timer is started on the first call to Ticks.
// If d == 0, the ticker internal timer is not started, and no ticks are
// dispatched.
func NewTimer(d time.Duration) TimeTicker {
	t := &timeTickerImpl{
		resetCh: make(chan time.Duration),
	}
	t.duration.Store(int64(d))
	return t
}

func (t *timeTickerImpl) Ticks() iter.Seq[time.Time] {
	defer t.Start()
	return t.tickerImpl.Ticks()
}

// Start the loop tick dispatcher loop, if it is not yet running. If called on a
// stopped, the ticks are restarted with the last non-zero period.
func (t *timeTickerImpl) Start() {
	if !t.running.Swap(true) {
		t.runWg.Add(1)
		go t.run()
	}
}

// Stop stops the timer and terminates consumers.
func (t *timeTickerImpl) Stop() {
	t.Reset(0)
	t.tickerImpl.Stop()
}

// Reset changes the period of the currently running and future ticks.
// If d == 0, the ticker timer will be stopped. If called on a stopped
// ticker with d != 0, the ticks are restarted.
func (t *timeTickerImpl) Reset(d time.Duration) {
	if d != 0 {
		// Do not store 0, so that [Start] starts normally.
		t.duration.Store(int64(d))
	}
	select {
	case t.resetCh <- d:
		if d == 0 {
			t.runWg.Wait()
		}
	default:
		t.Start()
	}
}

func (t *timeTickerImpl) run() {
	defer t.running.Store(false)
	defer t.runWg.Done()
	d := time.Duration(t.duration.Load())
	if d == 0 {
		return
	}
	t.Tick(time.Now())

	timer := time.NewTicker(d)
	defer timer.Stop()
	for {
		select {
		case tick, ok := <-timer.C:
			if !ok {
				return
			}
			t.Tick(tick)
		case d := <-t.resetCh:
			if d == 0 {
				return
			}
			timer.Reset(time.Duration(d))
		}
	}
}
