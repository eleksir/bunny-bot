package mylog

// Types only. The level constants, the name table and Level.name live in
// globals.go and level.go respectively, because a file named "types" must hold
// nothing that is not a type.

import (
	"io"
	stdlog "log"
	"os"
)

// Level is the severity of a log message. Levels are ordered, and a message is
// written when its own level is at or above the configured level.
type Level int

// stdLogger is the standard library logger doing the actual writing.
type stdLogger = stdlog.Logger

// Logger is the interface for modules that would otherwise depend on the
// standard library's *log.Logger. It is deliberately the smallest method set
// that *log.Logger satisfies, so widening a module's parameter type from
// *log.Logger to Logger keeps existing callers compiling.
type Logger interface {
	Output(calldepth int, s string) error
}

// adapter is the Logger handed to such modules. Messages arrive through
// Output and are logged at debug level into the already-open destination.
type adapter struct{}

// state is everything the logger needs to write a line. It is held behind a
// lock so that a concurrent Init cannot close a destination mid-write.
type state struct {
	level  Level
	dest   io.Writer
	slog   *stdLogger
	path   string
	owned  *os.File
	closed bool
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
