package ticker

import (
	"sync/atomic"
	"time"
)

type timeTickerImpl struct {
	tickerImpl[time.Time]
	resetCh  chan time.Duration
	duration atomic.Int64
}

var _ TimeTicker = (*timeTickerImpl)(nil)

func NewTimer(d time.Duration) *timeTickerImpl {
	t := &timeTickerImpl{
		tickerImpl: *New[time.Time](),
		resetCh:    make(chan time.Duration),
	}
	t.Reset(d)
	return t
}

func (t *timeTickerImpl) Start() {
	go t.run(time.Duration(t.duration.Load()))
}

// Stop stops the timer and terminates consumers.
func (t *timeTickerImpl) Stop() {
	t.Reset(0)
	t.tickerImpl.Stop()
}

// Reset changes the period of the currently running and future ticks.
func (t *timeTickerImpl) Reset(d time.Duration) {
	start := t.duration.Swap(int64(d)) == 0
	select {
	case t.resetCh <- d:
	default:
		if start && d != 0 {
			t.Start()
		}
	}
}

func (t *timeTickerImpl) run(d time.Duration) {
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
		case d, ok := <-t.resetCh:
			if !ok {
				return
			}
			if d == 0 {
				timer.Stop()
				return
			} else {
				timer.Reset(time.Duration(d))
			}
		}
	}
}
