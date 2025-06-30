package ticker

import (
	"sync/atomic"
	"testing"

	"github.com/parametalol/curry/assert"
)

func Test_consumer(t *testing.T) {
	t.Run("test send and ticks", func(t *testing.T) {
		c := newConsumer[int32]()
		i := atomic.Int32{}
		done := make(chan struct{})
		go func() {
			for x := range c.ticks() {
				i.Add(x)
			}
			close(done)
		}()
		c.send(1)
		c.send(10)
		c.send(100)
		c.close()
		<-done
		assert.That(t,
			assert.Equal(int32(111), i.Load()))
	})

	t.Run("close while sending", func(t *testing.T) {
		c := newConsumer[int]()
		done := make(chan struct{})
		go func() {
			done <- struct{}{}
			c.send(0)
			done <- struct{}{}
		}()
		<-done
		c.close()
		<-done
	})

	t.Run("send after done", func(t *testing.T) {
		c := newConsumer[int]()
		go c.send(0)
		for range c.ticks() {
			break
		}
		c.send(0)
		c.send(0)
	})
}
