package moon

// The worker registry: how a background goroutine joins the shutdown sequence,
// and how it leaves it. The state itself lives in globals.go with the rest of the
// package-level variables; this file holds only what a goroutine does with it.

// stopWorker unregisters a background goroutine that has finished.
func stopWorker() {
	runningWorkers.Add(-1)
	workers.Done()
}

// workerRunning reports whether any registered worker is still running.
func workerRunning() bool {
	return runningWorkers.Load() > 0
}

// startWorker registers a background goroutine with the shutdown sequence and
// reports whether it may run. It returns false once shutdown has begun, in
// which case the caller must return without starting: the stop signal is
// already closed and the teardown may already have closed what it would use.
func startWorker() bool {
	workersMu.Lock()
	defer workersMu.Unlock()

	if stopping {
		return false
	}

	workers.Add(1)
	runningWorkers.Add(1)

	return true
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
