package moon

import (
	"os"
	"os/signal"

	"bunny-bot/internal/mylog"
)

// InstallSignalHandler registers every signal the bot cares about and starts the
// handler goroutine. The signal sets differ per platform, so the knowledge lives
// behind terminationSignals and reopenSignals rather than leaking into main.
func InstallSignalHandler() {
	signal.Notify(SigChan, terminationSignals()...)
	signal.Notify(SigChan, reopenSignals()...)

	go SigHandler()
}

// SigHandler handles signals until one asks the process to quit.
//
// It stops the program; it does not end it. os.Exit here would skip main's
// deferred cleanup, and would report success whether or not the teardown actually
// worked — a service manager told "clean shutdown" over a log nobody can write to
// again is worse than one told to restart us. So the handler runs the teardown and
// returns; main decides the exit status from its own call to Shutdown.
func SigHandler() {
	for s := range SigChan {
		if !handleSignal(s) {
			continue
		}

		if err := shutdown(); err != nil {
			// The main log is already closed by now, so this is the last line we
			// get to write; main still exits non-zero over the same failure.
			mylog.Errorf("Shutdown did not complete cleanly: %s", err)
		}

		return
	}
}

// handleSignal applies one signal and reports whether it means "shut down".
//
// Reopen signals re-open the log file, which is what a log rotation tool needs
// after it has moved the file aside. Neither of them is expected on Windows,
// and an unrecognised signal is ignored rather than treated as a quit.
func handleSignal(s os.Signal) bool {
	if containsSignal(reopenSignals(), s) {
		mylog.Infof("Got %s, re-opening the log file", s)

		if err := mylog.ReOpen(); err != nil {
			// Reopening can legitimately fail — the destination may not be a
			// file at all, or the logger may already be closed. Logging is
			// best-effort here because the log is exactly what just broke.
			mylog.Errorf("Unable to re-open the log file on %s: %s", s, err)
		} else {
			mylog.Info("Log file re-opened")
		}

		return false
	}

	if !containsSignal(terminationSignals(), s) {
		// Make new iteration since we 've got signal we not interested in.
		return false
	}

	// The status has to be recorded here, before the teardown runs: main reads
	// it after Shutdown returns, and neither of them can reconstruct which
	// signal triggered all of this by then.
	recordExitCode(s)

	mylog.Infof("Got %s, quitting", s)

	return true
}

// containsSignal reports whether sig is in the set.
func containsSignal(set []os.Signal, sig os.Signal) bool {
	for _, s := range set {
		if s == sig {
			return true
		}
	}

	return false
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
