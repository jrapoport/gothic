package log

import "context"

type loggerKey struct{}

// WithContext adds a logger to the context.
func WithContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, log)
}

// FromContext returns a logger from the context.
func FromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerKey{}).(Logger); ok {
		return l
	}
	return NewStdLoggerWithLevel(PanicLevel)
}
