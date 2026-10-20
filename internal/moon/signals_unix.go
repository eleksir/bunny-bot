//go:build !windows

package moon

import (
	"os"
	"syscall"
)

// terminationSignals are the signals that ask the bot to shut down.
func terminationSignals() []os.Signal {
	return []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
}

// reopenSignals are the signals that ask the bot to re-open its log file.
//
// SIGHUP is the conventional daemon signal for "re-read your configuration and
// re-open your logs"; SIGUSR1 is the other common choice for log re-opening
// alone. Both are accepted so the bot fits either convention an operator's
// logrotate configuration expects.
func reopenSignals() []os.Signal {
	return []os.Signal{syscall.SIGHUP, syscall.SIGUSR1}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
