//go:build windows

package moon

import "os"

// exitCodeFor maps one terminating signal to the status a stopped process
// should carry. Windows has no signal numbers to follow, so it imitates the
// Unix convention: Ctrl+C is an interrupt, and an interrupt is 130.
func exitCodeFor(s os.Signal) int {
	if s == os.Interrupt {
		return signalExitBase + 2
	}

	return 0
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
