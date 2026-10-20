//go:build !windows

package moon

import (
	"os"
	"syscall"
)

// exitCodeFor maps one terminating signal to the status a stopped process
// should carry: 128 + the signal's number on Unix.
func exitCodeFor(s os.Signal) int {
	switch s {
	case syscall.SIGINT:
		return signalExitBase + 2

	case syscall.SIGTERM:
		return signalExitBase + 15

	case syscall.SIGQUIT:
		return signalExitBase + 3

	default:
		return 0
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
