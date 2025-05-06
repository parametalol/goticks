package goticks

type options struct {
	onStart    func() error
	onStop     func()
	stopTicker bool
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

func WithTickerStop() option {
	return func(o *options) {
		o.stopTicker = true
	}
}
