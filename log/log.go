package log

// Log is the same as the default standard logger from "log".
var Log = New()

// New returns a new log
func New() Logger {
	return NewStdLoggerWithLevel(PanicLevel)
}

// WithName returns a new named logger
func WithName(name string) Logger {
	return Log.WithName(name)
}

// Debug logs args when the logger level is debug.
func Debug(v ...interface{}) {
	Log.Debug(v...)
}

// Debugf formats args and logs the result when the logger level is debug.
func Debugf(format string, v ...interface{}) {
	Log.Debugf(format, v...)
}

// Info logs args when the logger level is info.
func Info(v ...interface{}) {
	Log.Info(v...)
}

// Infof formats args and logs the result when the logger level is info.
func Infof(format string, v ...interface{}) {
	Log.Infof(format, v...)
}

// Warn logs args when the logger level is warn.
func Warn(v ...interface{}) {
	Log.Warn(v...)
}

// Warnf formats args and logs the result when the logger level is warn.
func Warnf(format string, v ...interface{}) {
	Log.Warnf(format, v...)
}

// Error logs args when the logger level is error.
func Error(v ...interface{}) {
	Log.Error(v...)
}

// Errorf formats args and logs the result when the logger level is debug.
func Errorf(format string, v ...interface{}) {
	Log.Errorf(format, v...)
}

// Panic logs args on panic.
func Panic(v ...interface{}) {
	Log.Panic(v...)
}

// Panicf formats args and logs the result on panic.
func Panicf(format string, v ...interface{}) {
	Log.Panicf(format, v...)
}

// Fatal logs args when the error is fatal.
func Fatal(v ...interface{}) {
	Log.Fatal(v...)
}

// Fatalf formats args and logs the result when the error is fatal.
func Fatalf(format string, v ...interface{}) {
	Log.Fatalf(format, v...)
}

// UseFileOutput formats args and logs the result when the error is fatal.
func UseFileOutput(name string) Logger {
	return Log.UseFileOutput(name)
}
