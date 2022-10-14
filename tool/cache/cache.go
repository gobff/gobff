package cache

import (
	"time"
)

type (
	item[V any] struct {
		value      V
		expiration time.Time
	}
	cache[V any] struct {
		current  item[V]
		duration time.Duration
	}
	Cache[V any] interface {
		Set(value V)
		Get() (value V, found bool)
	}
)

func NewCache[V any](duration time.Duration) Cache[V] {
	return &cache[V]{
		duration: duration,
	}
}

func (c *cache[V]) Set(value V) {
	c.current = item[V]{
		value:      value,
		expiration: time.Now().Add(c.duration),
	}
}

func (c *cache[V]) Get() (V, bool) {
	if c.current.expiration.Before(time.Now()) {
		return *new(V), false
	}
	return c.current.value, true
}
