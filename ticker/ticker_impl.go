package ticker

import (
	"iter"
	"sync"
	"sync/atomic"
)

type tickerImpl[TickType any] struct {
	consumerID atomic.Int64
	consumers  sync.Map

	wg sync.WaitGroup
}

var _ Ticker[any] = (*tickerImpl[any])(nil)

func New[TickType any]() Ticker[TickType] {
	return &tickerImpl[TickType]{}
}

// Stop terminates consumers.
func (t *tickerImpl[TickType]) Stop() {
	t.forEach(func(id int64, consumer *consumer[TickType]) {
		t.consumers.Delete(id)
		consumer.close()
	})
}

// forEach executes f on every consumer.
func (t *tickerImpl[TickType]) forEach(f func(int64, *consumer[TickType])) {
	t.consumers.Range(func(key, value any) bool {
		f(key.(int64), value.(*consumer[TickType]))
		return true
	})
}

// Tick sends a tick to the consumers.
// It returns a [Waitable] on which the client may wait for the consumer to
// process the tick.
func (t *tickerImpl[TickType]) Tick(tick TickType) Waitable {
	tickWg := &sync.WaitGroup{}
	t.forEach(func(_ int64, consumer *consumer[TickType]) {
		tickWg.Add(1)
		t.wg.Add(1)
		go func() {
			consumer.send(tick)
			tickWg.Done()
			t.wg.Done()
		}()
	})
	return tickWg
}

// Ticks return a new iterator over the ticks.
func (t *tickerImpl[TickType]) Ticks() iter.Seq[TickType] {
	consumer := newConsumer[TickType]()
	t.consumers.Store(t.consumerID.Add(1), consumer)
	return consumer.ticks()
}

// Wait for the consumers to finish processing the current tick.
func (t *tickerImpl[TickType]) Wait() {
	t.wg.Wait()
}
