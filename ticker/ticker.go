package ticker

import (
	"iter"

	"time"
)

type Ticker[TickType any] interface {
	Ticks() iter.Seq[TickType]
	Stop()
	Wait()
}

type TickerWithTick[TickType any] interface {
	Ticker[TickType]
	Tick(TickType) interface{ Wait() }
}

type TimeTicker interface {
	Ticker[time.Time]
	Start()
	Reset(time.Duration)
}
