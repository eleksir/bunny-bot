package elog

import "time"

// nowUTC is the single clock every path and timestamp derives from.
//
// It is not time.Now at each call site: a log line stamped in local time beside a
// directory named for UTC is a bug report waiting to happen, and a file that rolls
// over at midnight local time while its lines are stamped UTC is worse.
func nowUTC() time.Time {
	return time.Now().UTC()
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
