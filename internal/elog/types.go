package elog

import (
	"sync"
	"time"
)

// chatLock serialises writes for one chat so that two goroutines cannot splice
// their lines into the same file.
type chatLock struct {
	mu sync.Mutex

	// refs counts goroutines holding or waiting for this entry. It is only
	// ever read or written while registryMu is held, which is what makes it a
	// usable proof of "nobody is holding this" during eviction.
	refs int

	lastUsed time.Time
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
