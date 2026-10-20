//go:build windows

package moon

import "os"

// terminationSignals are the signals that ask the bot to shut down.
func terminationSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}

// reopenSignals is empty on Windows.
//
// SIGHUP and SIGUSR1 are Unix signals that Windows never delivers — SIGUSR1 does
// not even exist in the syscall package there — so there is nothing to re-open
// the log for. The empty set keeps handleSignal platform-agnostic: a Windows
// service restarts to rotate its logs instead of signalling.
func reopenSignals() []os.Signal {
	return nil
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
