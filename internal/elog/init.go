package elog

import (
	"fmt"
	"os"
)

// Init switches chat logging on.
//
// An empty root disables it, and every Log* call becomes a no-op — chat logging
// is optional. A non-empty root must already exist: a path that does not is a
// typo, and silently logging nowhere is the worst possible response to one, so
// this returns an error and moon.ParseConfig treats it as fatal.
func Init(rootDir string) error {
	if rootDir == "" {
		setRoot("")

		return nil
	}

	info, err := os.Stat(rootDir)
	if err != nil {
		return fmt.Errorf("chat_logs directory %s is not usable: %w", rootDir, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: %s", errNotDirectory, rootDir)
	}

	setRoot(rootDir)

	return nil
}

// Close disables chat logging.
//
// Nothing is held open between lines, so there are no handles to close; this
// exists so shutdown has a symmetric counterpart to Init and so a misbehaving
// caller can stop writing.
func Close() error {
	setRoot("")

	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
