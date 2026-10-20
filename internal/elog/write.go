package elog

import (
	"os"
	"path/filepath"
	"time"

	"bunny-bot/internal/mylog"
)

// acquire returns the chat's lock with this goroutine holding a reference.
func acquire(chatID int64) *chatLock {
	registryMu.Lock()

	l, ok := locks[chatID]
	if !ok {
		l = &chatLock{}
		locks[chatID] = l
	}

	l.refs++
	l.lastUsed = time.Now()

	evictIdleLocked()

	registryMu.Unlock()

	l.mu.Lock()

	return l
}

// release drops this goroutine's reference and unlocks the chat.
func release(l *chatLock) {
	l.mu.Unlock()

	registryMu.Lock()

	l.refs--
	l.lastUsed = time.Now()

	registryMu.Unlock()
}

// evictIdleLocked drops registry entries nobody is holding and that have been
// idle for longer than chatLockIdleTimeout.
//
// Safety rests on refs: it is only ever changed while registryMu is held, so
// observing refs == 0 here is a proof that no goroutine holds the entry, not a
// guess. Evicting a held entry would let a second writer take a *different*
// mutex for the same chat, which is precisely the interleaving the lock exists
// to prevent. The caller must hold registryMu.
func evictIdleLocked() {
	cutoff := time.Now().Add(-chatLockIdleTimeout)

	for id, l := range locks {
		if l.refs == 0 && l.lastUsed.Before(cutoff) {
			delete(locks, id)
		}
	}
}

// appendLine writes one line to the chat's log for the day containing now.
//
// The whole protocol is set_mutex -> MkdirAll -> open|create -> write -> close,
// with unset_mutex deferred. Nothing is held open between lines.
func appendLine(chatID int64, line string) {
	// One read of both the root and whether logging is on, so a Close racing
	// this call can never leave us with an empty root and a relative path.
	dir, on := chatLogRoot()
	if !on {
		return
	}

	lock := acquire(chatID)

	defer release(lock)

	// One clock for both the path and the line: using Message.Date for one and
	// time.Now() for the other would let a line carry a date its own path
	// disagrees with.
	now := nowUTC()

	dayDir := filepath.Join(dir, itoa(chatID), now.Format("2006"), now.Format("01"))
	path := filepath.Join(dayDir, now.Format("02")+".txt")

	// MkdirAll creates every missing parent, so this covers the chat directory
	// and the year too. No need to pre-create anything.
	if err := os.MkdirAll(dayDir, logDirPerm); err != nil {
		reportFailure(err)

		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, logFilePerm)
	if err != nil {
		reportFailure(err)

		return
	}

	_, writeErr := f.WriteString(line + "\n")
	closeErr := f.Close()

	if writeErr != nil {
		reportFailure(writeErr)

		return
	}

	if closeErr != nil {
		reportFailure(closeErr)

		return
	}

	reportRecovered()
}

// reportFailure sends a write failure to the main log, but only the first one
// until a write succeeds again. A full disk would otherwise produce one error
// per chat event and drown the very log meant to explain it.
func reportFailure(err error) {
	stateMu.Lock()
	already := failing
	failing = true
	stateMu.Unlock()

	if already {
		return
	}

	mylog.Errorf("CHAT LOG WRITE FAILED, further failures suppressed until a write succeeds: %s", err)
}

// reportRecovered clears the suppression after a successful write.
func reportRecovered() {
	stateMu.Lock()
	was := failing
	failing = false
	stateMu.Unlock()

	if was {
		mylog.Info("Chat logging recovered")
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
