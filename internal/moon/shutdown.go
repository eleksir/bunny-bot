package moon

import (
	"bunny-bot/internal/elog"
	"bunny-bot/internal/mylog"
	"fmt"

	"github.com/akyoto/cache"
)

// closeCaches shuts down every cache sweeper goroutine.
//
// Safe only once the update listener and the greeting sweeper have stopped:
// cache.Close() reassigns the underlying sync.Map, which is a data race against
// any concurrent reader. shutdown() enforces that order.
func closeCaches() {
	closeOnce.Do(func() {
		for _, c := range []*cache.Cache{
			SquashedMembers,
			PendingGreetings.cache, GreetAnswers.cache,
		} {
			c.Close()
		}
	})
}

// Shutdown performs the ordered teardown, waits for it to finish, and reports
// whether every stage of it succeeded.
//
// main calls it after Telegram returns, for two reasons. The listener returning is
// not by itself a finished shutdown — a signal handler may be halfway through the
// sequence at that moment, and main returning would end the process underneath it,
// skipping stages and closing the log twice. And if the listener stopped for any
// other reason, nobody would run the sequence at all.
//
// It is safe to call from both goroutines: sync.Once runs the body once and blocks
// every other caller until that run completes, so both callers see the same result.
func Shutdown() error {
	return shutdown()
}

// shutdown stops everything in dependency order and waits for each stage.
//
// The order is not cosmetic. Every stage waits for the previous one to have
// finished, so no goroutine is ever using a resource that the next stage is
// closing:
//
//  1. the Telegram update listener — nothing else may add work;
//  2. the greeting sweeper — the last periodic user of the caches;
//  3. the caches — cache.Close() both signals the sweeper goroutines and
//     reassigns the underlying sync.Map, so it must run with nothing else
//     touching them;
//  4. chat logging, then 5. the main log — so the last log lines, which
//     describe stages 1 to 3, still have somewhere to go.
//
// Stages 1 and 2 share one stop signal and one WaitGroup, so a worker that has
// not started costs nothing and a worker that starts late is refused rather than
// racing the wait.
//
// The main log's own failure is reported ahead of the chat log's, because it happened
// last: a chat-log failure may no longer have been writable either.
func shutdown() error {
	var result error

	shutdownOnce.Do(func() {
		// Refuse new workers before signalling, so nothing can register between
		// the signal and the wait below.
		workersMu.Lock()
		stopping = true
		workersMu.Unlock()

		// Ask every loop to return, then wait for all of them. Nothing may be
		// using a cache, a log or the API once the next stages start closing.
		close(stopRequested)
		workers.Wait()

		closeCaches()

		// Chat logging first: it writes to the main log, so it must be closed
		// before the main log is.
		chatErr := elog.Close()

		// Since we shutting down, we have noting to do with log sync.
		logErr := mylog.Close()

		switch {
		case logErr != nil:
			result = fmt.Errorf("main log: %w", logErr)
		case chatErr != nil:
			result = fmt.Errorf("chat log: %w", chatErr)
		}
	})

	return result
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
