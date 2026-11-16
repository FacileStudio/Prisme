package unfurl

import (
	"sync"
	"time"
)

type entry struct {
	preview Preview
	expires time.Time
}

type cache struct {
	mu      sync.Mutex
	max     int
	ttl     time.Duration
	entries map[string]entry
}

func newCache(max int, ttl time.Duration) *cache {
	return &cache{max: max, ttl: ttl, entries: map[string]entry{}}
}

func (c *cache) get(key string) (Preview, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	found, ok := c.entries[key]
	if !ok || time.Now().After(found.expires) {
		delete(c.entries, key)
		return Preview{}, false
	}
	return found.preview, true
}

func (c *cache) put(key string, preview Preview) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= c.max {
		c.evict()
	}
	c.entries[key] = entry{preview: preview, expires: time.Now().Add(c.ttl)}
}

func (c *cache) evict() {
	now := time.Now()
	for key, found := range c.entries {
		if now.After(found.expires) {
			delete(c.entries, key)
		}
	}
	for key := range c.entries {
		if len(c.entries) < c.max {
			return
		}
		delete(c.entries, key)
	}
}
