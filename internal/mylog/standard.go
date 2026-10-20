package mylog

// The standard-library-compatible Logger: what a module receives instead of
// depending on *log.Logger itself, and where its messages end up.

import "strings"

// Standard returns the Logger for modules that would otherwise depend on the
// standard library's *log.Logger. Messages passed to it are logged at debug
// level into the already-open destination.
//
// Returning the interface rather than the concrete *adapter is the whole point:
// it is the type such modules declare, so hiding it would defeat the
// substitution it exists for.
func Standard() Logger { //nolint:ireturn // deliberate public interface
	return standard
}

// Output implements Logger, logging at debug level. The trailing newline of the
// standard library convention is trimmed so that mylog does not append a blank
// line to every message.
func (a *adapter) Output(_ int, s string) error {
	write(levelDebug, strings.TrimRight(s, "\r\n"))

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
