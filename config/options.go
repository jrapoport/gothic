package config

type configOptions struct {
	required bool
}

var defaultOptions = configOptions{
	required: true,
}

// Option is an interface to apply configOptions
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

// SkipRequired if set loading the config will not fail if required settings are missing.
func SkipRequired() Option {
	return newFuncOption(func(o *configOptions) {
		o.required = false
	})
}
