package ticker

import "iter"

type consumer[TickType any] struct {
	tickCh  chan TickType
	ackCh   chan struct{}
	closeCh chan struct{}
	doneCh  chan struct{}
}

func newConsumer[TickType any]() *consumer[TickType] {
	return &consumer[TickType]{
		make(chan TickType),
		make(chan struct{}),
		make(chan struct{}),
		make(chan struct{}),
	}
}

// send is the writer method that sends ticks to the consumer.
func (c *consumer[TickType]) send(tick TickType) {
	select {
	case <-c.doneCh:
	case <-c.closeCh:
		close(c.tickCh)
	case c.tickCh <- tick:
		<-c.ackCh
	}
}

// close is the writer method that closes the consumer.
func (c *consumer[TickType]) close() {
	close(c.closeCh)
}

// ticks returns an iterator that consumes all ticks and notifies the writer
// when the tick is processed.
func (c *consumer[TickType]) ticks() iter.Seq[TickType] {
	return func(f func(t TickType) bool) {
		defer close(c.ackCh)
		defer close(c.doneCh)
		for {
			select {
			case tick, ok := <-c.tickCh:
				if ok {
					ok = f(tick)
					c.ackCh <- struct{}{}
					if ok {
						continue
					}
				}
				return
			case <-c.closeCh:
				return
			}
		}
	}
}
