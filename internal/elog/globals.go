package elog

import (
	"errors"
	"sync"
	"time"
)

const (
	// chatLockIdleTimeout is how long a chat's lock entry may sit unused before it
	// is evicted. Without it the registry grows by one entry per chat ever seen.
	chatLockIdleTimeout = time.Hour

	// stampLayout is the per-line timestamp. GMT+0 with millisecond precision.
	stampLayout = "15:04:05.000"

	// logFilePerm is the mode for a newly created chat log file.
	logFilePerm = 0o644

	// logDirPerm is the mode for the chat, year and month directories.
	//
	// It is deliberately not logFilePerm. A directory needs the execute bit to
	// be traversable, and 0o644 has none, so creating the tree with the file mode
	// makes every subsequent open inside it fail with EACCES rather than ENOENT
	// — a failure reportFailure would swallow, leaving chat logging to silently
	// produce no files at all.
	//
	// The two live in one block on purpose: they are a pair that must never be
	// confused with each other, and keeping them adjacent is what makes confusing
	// them hard. They are here, in globals.go, because that is where every
	// package-level const belongs — the pair argument was weighed against rule 2
	// and rule 2 won, once, on the record.
	logDirPerm = 0o755
)

var (
	// errNotDirectory marks a chat_logs path that exists but is a regular file.
	errNotDirectory = errors.New("chat_logs path is not a directory")

	// root is the configured chat-logs directory. Empty means chat logging is
	// disabled and every Log* call is a no-op.
	root string

	// failing records that a write error has been reported and not yet cleared,
	// so a persistent failure (a full disk) logs once instead of once per event.
	failing bool

	// stateMu guards root and failing.
	stateMu sync.RWMutex

	// registryMu guards locks. It is never held while a chat lock is acquired,
	// or the whole logger would serialise.
	registryMu sync.Mutex

	// locks is the per-chat lock registry.
	locks = map[int64]*chatLock{}
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
