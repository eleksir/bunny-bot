package mylog

// The logger state: the one struct every write goes through, how a new one is
// built, and the two accessors that read and replace it under the package lock.
// The state itself lives in globals.go with the rest of the package-level
// variables; this file holds only what a caller does with it.

import (
	"io"
	stdlog "log"
	"os"
)

// newState builds a state whose stdlib logger writes to dest.
func newState(lvl Level, dest io.Writer, path string, owned *os.File, closed bool) *state {
	return &state{
		level:  lvl,
		dest:   dest,
		slog:   newStdLogger(dest),
		path:   path,
		owned:  owned,
		closed: closed,
	}
}

// newStdLogger wraps a writer in a standard library logger. Flags are zeroed
// because the line format is rendered by mylog itself.
func newStdLogger(dest io.Writer) *stdLogger {
	l := stdlog.New(dest, "", 0)
	l.SetFlags(0)

	return l
}

// getState returns the active logger state.
func getState() *state {
	lock.RLock()
	defer lock.RUnlock()

	return current
}

// setState replaces the active state wholesale. Tests use it to restore the
// documented before-Init default.
func setState(s *state) {
	lock.Lock()
	defer lock.Unlock()

	current = s
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
