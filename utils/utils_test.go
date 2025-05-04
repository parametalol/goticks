package utils

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSeqIgnoreErr(t *testing.T) {
	i := 2
	inc := func() {
		i++
	}
	mul := func(_ context.Context, _ any) error {
		i *= 2
		return errors.New("error")
	}
	assert.NoError(t, Seq(Adapt[any](inc), IgnoreErr[any](mul))(context.Background(), 0))
	assert.Equal(t, 6, i)

	assert.Error(t, Seq(mul, Adapt[any](inc))(context.Background(), 0))
	assert.Equal(t, 12, i)
}

type arr []string

func (a *arr) Write(data []byte) (int, error) {
	*a = append(*a, string(data))
	return len(data), nil
}

func TestWithLog(t *testing.T) {
	var a = &arr{}
	err := WithLog[any](a, a, "test", func() error {
		return errors.New("test")
	})(context.Background(), nil)

	assert.Error(t, err)
	assert.Equal(t, arr{
		"Calling test\n",
		"Execution of test failed with error: test\n",
	}, *a)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err = WithLog[any](a, a, "test", func(context.Context) {})(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, arr{
		"Calling test\n",
		"Execution cancelled for test\n",
	}, (*a)[2:])
}

func TestWithTimeout(t *testing.T) {
	var deadline time.Time
	var ok bool
	now := time.Now()
	err := WithTimeout[any](0, func(ctx context.Context) error {
		deadline, ok = ctx.Deadline()
		return ctx.Err()
	})(context.Background(), 0)
	assert.ErrorIs(t, err, context.DeadlineExceeded)
	assert.True(t, ok)
	assert.True(t, time.Since(now) >= time.Since(deadline))
}

func TestNoOverlap(t *testing.T) {
	var i atomic.Int32
	testCh := make(chan bool)
	task := func() {
		i.Add(1)
		testCh <- true
		testCh <- true
	}
	fn := NoOverlap[any](task)
	go func() {
		_ = fn(context.Background(), 0)
	}()
	<-testCh
	_ = fn(context.Background(), 0)
	_ = fn(context.Background(), 0)
	_ = fn(context.Background(), 0)
	<-testCh
	assert.Equal(t, int32(1), i.Load())
}

func TestWithRetry(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		var i int
		task := func() error {
			i++
			return errors.New("test")
		}
		err := WithRetry[any](SimpleRetryPolicy(3), task)(context.Background(), 0)
		assert.Error(t, err)
		assert.Equal(t, 3, i)
	})
	t.Run("without error", func(t *testing.T) {
		var i int
		task := func() {
			i++
		}
		err := WithRetry[any](SimpleRetryPolicy(3), task)(context.Background(), 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, i)
	})
	t.Run("with cancelled context", func(t *testing.T) {
		var i int
		task := func() {
			i++
		}
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		err := WithRetry[any](SimpleRetryPolicy(3), task)(ctx, 0)
		assert.NoError(t, err)
		assert.Equal(t, 1, i)
	})
	t.Run("with exponential backoff", func(t *testing.T) {
		var i int
		task := func() error {
			i++
			return errors.New("test")
		}
		err := WithRetry[any](ExponentialBackoffPolicy(3, time.Millisecond), task)(context.Background(), 0)
		assert.Error(t, err)
		assert.Equal(t, 3, i)
	})
}
