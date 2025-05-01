package ticker

import "iter"

type tack[TickType any] struct {
	tick  TickType
	ackCh chan struct{}
}

// consumer wraps a tick channel and synchronously acknowledges the tick
// processing.
type consumer[TickType any] struct {
	tickCh  chan tack[TickType]
	closeCh chan struct{}
	doneCh  chan struct{}
}

func newConsumer[TickType any]() *consumer[TickType] {
	return &consumer[TickType]{
		tickCh:  make(chan tack[TickType]),
		closeCh: make(chan struct{}),
		doneCh:  make(chan struct{}),
	}
}

// send is the writer method that sends ticks to the consumer.
func (c *consumer[TickType]) send(tick TickType) {
	tack := tack[TickType]{tick, make(chan struct{})}
	select {
	case <-c.doneCh:
	case <-c.closeCh:
		close(c.tickCh)
	case c.tickCh <- tack:
		<-tack.ackCh
	}
}

// close is the writer method that closes the consumer.
// The closed consumer won't receive more ticks, and cannot be reopened.
func (c *consumer[TickType]) close() {
	close(c.closeCh)
}

// ticks returns an iterator that consumes all ticks and notifies the writer
// when the tick is processed.
func (c *consumer[TickType]) ticks() iter.Seq[TickType] {
	return func(yield func(t TickType) bool) {
		defer close(c.doneCh)
		for {
			select {
			case tickAck, ok := <-c.tickCh:
				if !ok {
					return
				}
				ok = yield(tickAck.tick)
				close(tickAck.ackCh)
				if !ok {
					return
				}
			case <-c.closeCh:
				return
			}
		}
	}
}
