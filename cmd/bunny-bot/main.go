package main

import (
	"fmt"
	"os"
	"path/filepath"

	"bunny-bot/internal/moon"
	"bunny-bot/internal/mylog"
)

// main decides the process's fate and nothing else: run the bot, and decide the
// exit status from how it ended. An error leaves 1; a signal-driven stop leaves
// 128 + the signal's number (130 for Ctrl+C on this platform) so a supervisor can
// tell a signal stop from a clean one, which is what D33 fixed; a stop nobody
// signalled leaves 0.
//
// os.Exit is called here and nowhere else, deliberately. A worker goroutine cannot
// set an exit status without ending the process, which skips every deferred
// cleanup on the way out — including the one that closes the log we would be
// complaining about.
func main() {
	if err := run(); err != nil {
		// The logs may already be closed, so stderr is the one channel left that
		// is certain to be attached.
		fmt.Fprintf(os.Stderr, "bunny-bot: %v\n", err)

		os.Exit(1)
	}

	os.Exit(moon.ExitCode())
}

// run is the whole lifecycle, so that every deferred cleanup belongs to a function
// that returns instead of one that exits.
//
// Its single defer covers the bot's entire life: once the log is open, nothing
// after that point returns early.
func run() error {
	// Main config.
	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("unable to get current executable path: %w", err)
	}

	configJSONPath := fmt.Sprintf("%s/data/config.json", filepath.Dir(executablePath))

	moon.Cfg, err = moon.ParseConfig(configJSONPath)
	if err != nil {
		return fmt.Errorf("unable to read the config: %w", err)
	}

	// Open the log. ParseConfig has already coerced an unknown level to "info",
	// so an error here means the file itself could not be opened — which is
	// worth surfacing loudly rather than continuing with no record of anything.
	if err := mylog.Init(moon.Cfg.LogLevel, moon.Cfg.LogFile); err != nil {
		return fmt.Errorf("unable to open the log: %w", err)
	}

	defer func() { _ = mylog.Close() }()

	mylog.Info("Config read looks good")

	// Самое время поставить траппер сигналов.
	moon.InstallSignalHandler()

	mylog.Info("Signal handler installed")

	mylog.Info("Starting main loop")

	// Telegram assigns moon.Tg on this goroutine; the sweeper reads it, so the
	// client must be constructed before any goroutine that can call the API.
	moon.Telegram(moon.Cfg)

	// The update listener has returned, but that is not the same as a finished
	// shutdown: a signal handler may be mid-teardown right now. Shutdown blocks
	// until the ordered sequence has completed, whichever goroutine started it,
	// so the process cannot exit with the log half-closed — and its error is what
	// decides the exit status.
	return moon.Shutdown()
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
