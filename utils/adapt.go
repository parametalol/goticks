package utils

import (
	"context"
	"time"
)

// generic type alias requires GOEXPERIMENT=aliastypeparams
// type normalizedFunc[TickType any] = func(context.Context, TickType) error

type Func[TickType any] interface {
	~func(context.Context, TickType) error |
		~func() | ~func() error | ~func(context.Context) | ~func(context.Context) error |
		~func(TickType) | ~func(TickType) error | ~func(context.Context, TickType)
}

// Adapt the function to a signature that takes a context and returns an error.
//
// Example:
//
//	f := func(){}
//	err := Adapt[time.Time](f)(context.Background, time.Now())
func Adapt[TickType any, Fn Func[TickType]](f Fn) func(context.Context, TickType) error {
	switch t := any(f).(type) {
	case func(context.Context, TickType) error:
		return t
	case func():
		return func(_ context.Context, _ TickType) error {
			t()
			return nil
		}
	case func() error:
		return func(_ context.Context, _ TickType) error {
			return t()
		}
	case func(context.Context):
		return func(ctx context.Context, _ TickType) error {
			t(ctx)
			return nil
		}
	case func(context.Context) error:
		return func(ctx context.Context, _ TickType) error {
			return t(ctx)
		}
	case func(TickType):
		return func(_ context.Context, tick TickType) error {
			t(tick)
			return nil
		}
	case func(TickType) error:
		return func(_ context.Context, tick TickType) error {
			return t(tick)
		}
	case func(context.Context, TickType):
		return func(ctx context.Context, tick TickType) error {
			t(ctx, tick)
			return nil
		}
	}
	panic("unsupported function signature")
}

func AdaptT[Fn Func[time.Time]](f Fn) func(context.Context, time.Time) error {
	return Adapt[time.Time](f)
}
