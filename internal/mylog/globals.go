package mylog

import (
	"errors"
	"io"
	"os"
	"sync"
)

const (
	// levelTrace is the most verbose level, below DEBUG.
	levelTrace Level = iota
	levelDebug
	levelInfo
	levelWarn
	levelError
	levelFatal
)

const (
	// stampLayout is the timestamp layout. Sub-second precision is kept on
	// purpose: the greeting sweeper reasons over 15-minute windows measured in
	// microseconds, and the previous zap configuration truncated to whole
	// seconds, which made window-edge debugging needlessly hard.
	stampLayout = "2006-01-02 15:04:05.000"

	// calldepth is reported to the stdlib logger. mylog never prints caller
	// information, so the value only has to be consistent.
	calldepth = 2

	// filePerm is the mode for a newly created log file.
	filePerm = 0o644
)

var (
	// errClosed is returned when ReOpen is called on a closed logger.
	errClosed = errors.New("cannot reopen a closed logger")

	// errUnknownLevel is returned for a log level name mylog does not know.
	errUnknownLevel = errors.New("unknown log level")

	// levelNames renders each level as the uppercase token written to the log.
	levelNames = map[Level]string{
		levelTrace: "TRACE",
		levelDebug: "DEBUG",
		levelInfo:  "INFO",
		levelWarn:  "WARN",
		levelError: "ERROR",
		levelFatal: "FATAL",
	}

	// exit is os.Exit, indirected so tests can observe the fatal path.
	exit = os.Exit

	// lock guards the package-global logger state. It is held for reading across
	// the actual write, so that a concurrent Init cannot close the destination
	// between the level check and the write.
	lock sync.RWMutex

	// current is the active logger. It is deliberately never nil: it is created
	// working at package initialisation, writing to stderr at INFO. A nil logger
	// would be the defect M2 in this codebase, and a valid default removes the
	// state entirely rather than guarding against it.
	current = newState(levelInfo, io.Writer(os.Stderr), "", nil, false)

	// standard is the singleton Logger handed to modules that expect a
	// *log.Logger. It is stateless; the package-global current does the work.
	standard Logger = &adapter{}
)

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
