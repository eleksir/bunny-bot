package mylog

// The level vocabulary: what a level is called, how a name in free-form config
// text becomes a Level, and how the configured threshold is changed without
// reopening the destination.

import (
	"fmt"
	"strings"
)

// ParseLevel maps a level name to a Level. Parsing is case-insensitive and
// tolerates surrounding whitespace, because the level comes from free-form user
// text in the config file. An unknown name is an error: Init and SetLevel are
// programming-level calls where a typo should be loud.
func ParseLevel(s string) (Level, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "TRACE":
		return levelTrace, nil
	case "DEBUG":
		return levelDebug, nil
	case "INFO":
		return levelInfo, nil
	case "WARN":
		return levelWarn, nil
	case "ERROR":
		return levelError, nil
	case "FATAL":
		return levelFatal, nil
	default:
		return levelInfo, fmt.Errorf("%w %q, want one of trace, debug, info, warn, error, fatal", errUnknownLevel, s)
	}
}

// SetLevel changes the level without reopening the destination.
func SetLevel(loglevel string) error {
	lvl, err := ParseLevel(loglevel)
	if err != nil {
		return err
	}

	lock.Lock()
	defer lock.Unlock()

	current.level = lvl

	return nil
}

// name returns the uppercase token for the level.
func (l Level) name() string {
	if n, ok := levelNames[l]; ok {
		return n
	}

	return "UNKNOWN"
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
