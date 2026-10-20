package moon

import (
	"time"

	"github.com/akyoto/cache"
)

// newEventCache creates an eventCache whose expired entries are swept every
// cleaningInterval.
func newEventCache[T any](cleaningInterval time.Duration) *eventCache[T] {
	return &eventCache[T]{cache: cache.New(cleaningInterval)}
}

// Set stores value under the microsecond key with the given TTL.
func (m *eventCache[T]) Set(micros int64, value T, ttl time.Duration) {
	m.cache.Set(micros, value, ttl)
}

// Snapshot returns a typed copy of all non-expired entries keyed by their
// microsecond timestamp. Callers can range over the returned map without
// performing any interface{} conversions.
func (m *eventCache[T]) Snapshot() map[int64]T {
	out := make(map[int64]T)

	m.cache.Range(
		func(k, v any) bool {
			key, ok := k.(int64)
			if !ok {
				return true
			}

			value, ok := v.(T)
			if !ok {
				return true
			}

			out[key] = value

			return true
		},
	)

	return out
}

// Delete removes the entry stored under the given microsecond key.
func (m *eventCache[T]) Delete(micros int64) {
	m.cache.Delete(micros)
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
