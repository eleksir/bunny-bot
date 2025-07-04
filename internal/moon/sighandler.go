package moon

import (
	"os"
	"syscall"
)

// SigHandler.handles signals. Mainly closes persistent dbs and quits for good.
func SigHandler() {
	for {
		var s = <-SigChan
		switch s {
		case syscall.SIGINT:
			Log.Infoln("Got SIGINT, quitting")
		case syscall.SIGTERM:
			Log.Infoln("Got SIGTERM, quitting")
		case syscall.SIGQUIT:
			Log.Infoln("Got SIGQUIT, quitting")

		// Make new iteration since we 've got signal we not interested in.
		default:
			continue
		}

		if len(DB) > 0 {
			Log.Debug("Closing runtime bot settings db")

			for _, db := range DB {
				_ = db.Close()
			}
		}

		// Since we shutting down, we have noting to do with log sync.
		_ = Log.Sync()

		os.Exit(0)
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
