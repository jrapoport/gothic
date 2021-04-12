package config

type configOptions struct {
	required bool
}

var defaultOptions = configOptions{
	required: true,
}

type Option interface {
	apply(*configOptions)
}

type funcOption struct {
	f func(*configOptions)
}

func newFuncOption(f func(*configOptions)) *funcOption {
	return &funcOption{
		f: f,
	}
}
func (fdo *funcOption) apply(do *configOptions) {
	fdo.f(do)
}

func SkipRequired() Option {
	return newFuncOption(func(o *configOptions) {
		o.required = false
	})
}
