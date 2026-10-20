package moon

import (
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/akyoto/cache"
)

const (
	// greetCacheTTL is how long a join event or a chat message is retained. It
	// is deliberately larger than greetReplyWindow so that a sweep occurring at
	// the far end of the window still has the data it needs to decide.
	greetCacheTTL = 17 * time.Minute

	// greetReplyWindow is how long a newcomer has to answer the greeting before
	// being treated as a spammer.
	greetReplyWindow = 15 * time.Minute

	// greetSweepInterval is how often the sweeper runs.
	greetSweepInterval = 1 * time.Minute
)

// signalExitBase is the offset a signal's exit status is derived from. The Unix
// convention is 128 + n, where n is the signal number, so SIGINT is 130, SIGTERM
// is 143 and SIGQUIT is 131 — what a shell reports for a process killed by a
// signal, which is exactly what these statuses imitate.
const signalExitBase = 128

// The outcomes of a squash attempt. Names state what happened, not how it was
// detected, so callers never have to reinterpret them.
const (
	// banApplied means Telegram accepted the ban.
	banApplied banOutcome = iota

	// banSuppressed means we did not call the API at all: the debounce already
	// holds this user, i.e. we banned them moments ago and they are already out.
	banSuppressed

	// banFailed means we called the API and it refused. The user is still in the
	// chat and was not retried.
	banFailed
)

// The outcomes of retracting a greeting a banned member never answered. The three
// cases are distinct because "nothing to delete" and "tried and failed" call for
// different log lines, and only the second is worth an operator's attention.
const (
	// nothingToRetract means no remembered greeting id was greater than zero, so
	// the API was never called.
	nothingToRetract retractOutcome = iota

	// greetingRetracted means every remembered greeting was deleted.
	greetingRetracted

	// greetingRetractFailed means at least one deletion was refused by Telegram.
	// The others still went: a failure never aborts the rest.
	greetingRetractFailed
)

var (
	// htmlReplacer escapes the three characters HTML requires, all three pairs at
	// once so no ordering can be got wrong: "&" is never re-escaped by the "<" or
	// ">" rules. Created once because a Replacer is immutable and safe to share.
	htmlReplacer = strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;")

	Cfg     MyConfig
	SigChan = make(chan os.Signal, 1)
	Tg      tgClient

	// allowedUpdates is the set of update types Telegram is asked to deliver.
	//
	// It lives here, not inline in Telegram(), because it must stay in step with
	// the branches of telegramMsgParser: a type missing from this list makes its
	// branch unreachable, and nothing at compile time or in the test suite says
	// so. TestEveryDispatchedUpdateTypeIsRequested exists to catch exactly that.
	//
	// chat_join_request is deliberately absent. The community is open and join
	// requests are not moderated; Telegram auto-approves them when the type is
	// not requested, which is the intended behaviour.
	allowedUpdates = []echotron.UpdateType{
		"message",
		"edited_message",
		"chat_member",
	}

	// SquashedMembers holds members recently banned by us, used as a ban debounce.
	SquashedMembers = cache.New(1 * time.Minute)

	// PendingGreetings holds recorded joins keyed by join timestamp in
	// microseconds. Used by the greeting sweeper to find newcomers who never
	// answered, and to retract the greetings we asked them.
	PendingGreetings = newEventCache[pendingGreeting](1 * time.Minute)

	// GreetAnswers holds text messages keyed by message timestamp in
	// microseconds. Used by the greeting sweeper to find proof that a newcomer
	// answered the greeting.
	GreetAnswers = newEventCache[*echotron.Message](1 * time.Minute)

	// closeOnce guards closeCaches: cache.Cache.Close() sends on an unbuffered
	// channel and therefore must be called exactly once per cache; calling it
	// twice deadlocks. Each cache is guarded so a failure partway through cannot
	// turn into a hang on the second attempt.
	closeOnce sync.Once

	// stopRequested is closed once, by shutdown, to ask every goroutine that
	// loops to return. One signal for all of them, rather than one per worker:
	// several workers can come and go (and tests start them deliberately), and
	// a per-worker channel would have to be closed exactly once each, which is a
	// panic waiting to happen.
	stopRequested = make(chan struct{})

	// workers counts the running background goroutines — the update listener and
	// the greeting sweeper. shutdown waits on it, so no worker can be touching a
	// cache, the logs or the API while the teardown closes them.
	workers sync.WaitGroup

	// workersMu guards workers and stopping. A worker must register before
	// shutdown starts waiting, and a goroutine that tries to start afterwards is
	// refused instead of racing WaitGroup.Add against Wait.
	workersMu sync.Mutex
	stopping  bool

	// runningWorkers counts registered, not-yet-finished workers, so a caller
	// can tell "nothing started yet" from "started and still going".
	runningWorkers atomic.Int32

	// exitCode carries the status the process should exit with on a clean
	// return from run(): zero while the bot is running and sleeping normally,
	// 128 + signal number once a termination signal has begun the shutdown.
	// It is an atomic because the handler goroutine writes it and main reads
	// it, and plain avoids any question of which goroutine wins.
	exitCode atomic.Int32

	// shutdownOnce guards shutdown so two termination signals arriving together
	// cannot run the sequence twice.
	shutdownOnce sync.Once
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
