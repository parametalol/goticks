package ticker

import (
	"iter"

	"time"
)

type Tickable[TickType any] interface {
	Ticks() iter.Seq[TickType]
	Tick(TickType) Waitable
}

type Startable interface {
	Start()
}

type Stoppable interface {
	Stop()
}

type Restartable interface {
	Startable
	Stoppable
}

type Waitable interface {
	Wait()
}

type Ticker[TickType any] interface {
	Tickable[TickType]
	Stoppable
	Waitable
}

type TimeTicker interface {
	Tickable[time.Time]
	Restartable
	Waitable
	Reset(time.Duration)
}
