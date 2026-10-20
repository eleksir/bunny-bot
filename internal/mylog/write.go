package mylog

// The one-line emit path: decide whether the line is wanted, render it, write
// it. Everything else in the package exists to reach these two functions.

import "time"

// write renders and writes one line if the level passes the filter and the
// logger is open. The read lock is held across the write so that a concurrent
// Init cannot close the destination underneath us.
func write(lvl Level, msg string) {
	lock.RLock()
	defer lock.RUnlock()

	if current.closed || lvl < current.level {
		return
	}

	_ = current.slog.Output(calldepth, render(lvl, msg))
}

// render builds one log line: a GMT+0 timestamp, the uppercase level, and the
// message verbatim, separated by exactly one space each.
func render(lvl Level, msg string) string {
	return time.Now().UTC().Format(stampLayout) + " " + lvl.name() + " " + msg
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
