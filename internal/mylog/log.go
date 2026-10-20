package mylog

// The destination lifecycle: opening the log, re-opening it for rotation, and
// closing it for good. The level vocabulary lives in level.go, the state in
// state.go, and the act of writing a line in write.go.

import (
	"fmt"
	"io"
	stdlog "log"
	"os"
)

// Init opens the log destination and sets the level.
//
// An empty filepath selects stderr. A non-empty path is opened appending, so
// an existing log is never truncated.
//
// Init is re-entrant: any previously open file is closed first, so repeated
// calls do not leak descriptors. If the new destination cannot be opened, Init
// returns an error and leaves the previous destination and level untouched —
// it deliberately does not fall back to stderr, because a silent fallback
// hides a broken log path from the operator.
func Init(loglevel string, filepath string) error {
	lvl, err := ParseLevel(loglevel)
	if err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	dest, owned := io.Writer(os.Stderr), (*os.File)(nil)

	if filepath != "" {
		// The path comes from the operator's own config file, not from a
		// network peer, and opening it is this function's entire purpose.
		file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, filePerm) //nolint: gosec // G703
		if err != nil {
			return fmt.Errorf("unable to open log file %s: %w", filepath, err)
		}

		dest, owned = file, file
	}

	swapLocked(lvl, dest, filepath, owned)

	// Keep the standard library's default logger pointed at the same
	// destination, so third-party stdlib logging does not bypass our file.
	stdlog.SetFlags(0)
	stdlog.SetOutput(dest)

	return nil
}

// ReOpen closes and re-opens the current log file, keeping its level. It exists
// for external rotation: the rotation tool moves the file aside, and ReOpen
// makes the logger pick up the new one.
//
// ReOpen on a closed logger is an error rather than a silent reopen, and on a
// stderr destination it is a no-op, since there is no file to re-open.
func ReOpen() error {
	lock.RLock()
	closed, path, lvl := current.closed, current.path, current.level
	lock.RUnlock()

	switch {
	case closed:
		return errClosed
	case path == "":
		return nil
	default:
		return Init(lvl.name(), path)
	}
}

// Close closes the log destination. It is idempotent, and it is safe when the
// destination is stderr — stderr is not ours to close.
//
// After Close, logging calls are silently dropped rather than panicking or
// resurrecting the file: Close is terminal.
func Close() error {
	lock.Lock()
	defer lock.Unlock()

	if current.closed {
		return nil
	}

	current.closed = true

	if current.owned == nil {
		return nil
	}

	// The stdlib Logger has no Close — it does not own its writer — so the file
	// we opened ourselves is closed directly.
	file := current.owned
	current.owned = nil

	if file == nil {
		return nil
	}

	return file.Close()
}

// swapLocked installs a new state, closing any file we previously owned.
// The caller must hold the write lock.
func swapLocked(lvl Level, dest io.Writer, path string, owned *os.File) {
	if current.owned != nil {
		_ = current.owned.Close()
	}

	current = newState(lvl, dest, path, owned, false)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
