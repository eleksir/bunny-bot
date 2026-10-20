package elog

// The value formatters: small conversions shared by the write path (which builds
// path components) and the renderer (which builds inline fields). They are shared
// by two topics, so they get a file of their own rather than living in one of them
// and making the other look in the wrong place.

import (
	"strconv"
	"strings"
)

// field renders "key=value", or empty when the value is empty. Empty fields are
// omitted: a bare "username=@" is worse than nothing.
//
// It lives here rather than in render.go because sizeField, a formatter, is one of
// its callers — and a helper two topics use belongs to neither.
func field(key string, value string) string {
	if value == "" {
		return ""
	}

	return key + "=" + value
}

// joinFields drops empty fragments and joins the rest with single spaces.
func joinFields(parts []string) string {
	kept := make([]string, 0, len(parts))

	for _, p := range parts {
		if p != "" {
			kept = append(kept, p)
		}
	}

	return strings.Join(kept, " ")
}

// itoa renders an int64 for a path component or an inline field.
func itoa(v int64) string {
	return strconv.FormatInt(v, 10)
}

// itoaInt renders an int for a rendered field.
func itoaInt(v int) string {
	return strconv.Itoa(v)
}

// itoaFloat renders a float for a rendered field.
func itoaFloat(v float64) string {
	return strconv.FormatFloat(v, 'g', -1, 64)
}

// atUsername renders a username with the conventional @ prefix, or empty.
func atUsername(name string) string {
	if name == "" {
		return ""
	}

	return "@" + name
}

// seconds renders a duration in seconds with its unit.
func seconds(n int) string {
	return itoaInt(n) + "s"
}

// sizeField renders a byte count, omitted when unknown. A size of 0 means
// Telegram did not tell us, not that the file was empty.
func sizeField(n int64) string {
	if n <= 0 {
		return ""
	}

	return field("size", itoa(n))
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
