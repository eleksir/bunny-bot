package mylog

import (
	"fmt"
)

// Trace logs at the most verbose level, below DEBUG.
func Trace(args ...any) { write(levelTrace, fmt.Sprint(args...)) }

// Tracef logs a formatted message at the most verbose level, below DEBUG.
func Tracef(format string, args ...any) { write(levelTrace, fmt.Sprintf(format, args...)) }

// Debug logs at debug level.
func Debug(args ...any) { write(levelDebug, fmt.Sprint(args...)) }

// Debugf logs a formatted message at debug level.
func Debugf(format string, args ...any) { write(levelDebug, fmt.Sprintf(format, args...)) }

// Info logs at info level.
func Info(args ...any) { write(levelInfo, fmt.Sprint(args...)) }

// Infof logs a formatted message at info level.
func Infof(format string, args ...any) { write(levelInfo, fmt.Sprintf(format, args...)) }

// Warn logs at warn level.
func Warn(args ...any) { write(levelWarn, fmt.Sprint(args...)) }

// Warnf logs a formatted message at warn level.
func Warnf(format string, args ...any) { write(levelWarn, fmt.Sprintf(format, args...)) }

// Error logs at error level.
func Error(args ...any) { write(levelError, fmt.Sprint(args...)) }

// Errorf logs a formatted message at error level.
func Errorf(format string, args ...any) { write(levelError, fmt.Sprintf(format, args...)) }

// Fatal logs at fatal level, then exits the process with status 1.
//
// The destination is closed before exiting. os.Exit does not run deferred
// functions, so without this an explicit close a buffered tail could be lost —
// which is precisely the case where the message matters most.
func Fatal(args ...any) {
	write(levelFatal, fmt.Sprint(args...))

	exitAfterClose(1)
}

// Fatalf logs a formatted message at fatal level, then exits the process with
// status 1. The destination is closed first; see Fatal.
func Fatalf(format string, args ...any) {
	write(levelFatal, fmt.Sprintf(format, args...))

	exitAfterClose(1)
}

// exitAfterClose closes the logger and then exits, in that order.
func exitAfterClose(code int) {
	_ = Close()

	exit(code)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
