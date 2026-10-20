package moon

import "os"

// ExitCode reports the status the process should leave with on a clean return
// from run(). It is zero while the bot is running normally, and 128 + the
// signal's number once a termination signal has begun the shutdown — the Unix
// convention a supervisor can use to tell "stopped with a signal" from "stopped
// cleanly", which D33 left permanently 0 before it.
//
// It never reports an error status: a failed shutdown still returns an error
// from run(), and main exits 1 over that, the same as before.
func ExitCode() int {
	return int(exitCode.Load())
}

// recordExitCode stores the status the given terminating signal should leave.
// Only the first termination signal matters — the handler goroutine returns
// after one — so a later write cannot pass either.
func recordExitCode(s os.Signal) {
	// The value is one of a fixed table of 128+n, never a caller-supplied int.
	exitCode.Store(int32(exitCodeFor(s))) //nolint:gosec // G115: exitCodeFor returns only 130, 131, 143 or 0.
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
