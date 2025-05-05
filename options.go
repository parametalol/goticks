package goticks

type options struct {
	onStart func() error
	onStop  func()
}

type option func(*options)

func WithOnStart(f func() error) option {
	return func(o *options) {
		o.onStart = f
	}
}

func WithOnStop(f func()) option {
	return func(o *options) {
		o.onStop = f
	}
}
