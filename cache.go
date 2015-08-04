// Package cache implements a (very) simple cache
// mechanism, suitable for single-process applications
package cache

import (
	"sync"
	"time"
)

var (
	defaultExpireTime          = time.Minute * 5
	defaultExpiryCheckInterval = time.Second * 10
)

// A Cache holds a list of Entries, which it occasionally purges
type Cache struct {
	mutex                  *sync.RWMutex
	entryExpirationDefault time.Duration
	expiryCheckInterval    time.Duration
	entries                map[string]Entry
}

// cleanup removes stale entries from the cache
func (c *Cache) cleanup() {
	c.mutex.Lock()
	for key, item := range c.entries {
		if item.expired() {
			delete(c.entries, key)
		}
	}
	c.mutex.Unlock()
}

// New returns a fresh Cache
func New(entryExpirationDefault, expiryCheckInterval *time.Duration) *Cache {
	// Set default entry expiration time if none provided
	if entryExpirationDefault == nil {
		entryExpirationDefault = &defaultExpireTime
	}

	// Set default expiry check interval if none provided
	if expiryCheckInterval == nil {
		expiryCheckInterval = &defaultExpiryCheckInterval
	}

	cache := &Cache{
		mutex: &sync.RWMutex{},
		entryExpirationDefault: *entryExpirationDefault,
		expiryCheckInterval:    *expiryCheckInterval,
		entries:                make(map[string]Entry),
	}

	// Start recurring clean-up loop
	ticker := time.Tick(cache.expiryCheckInterval)
	go func() {
		for {
			<-ticker
			cache.cleanup()
		}
	}()

	return cache
}

// Set a value for a key
func (c Cache) Set(key string, value string) {
	c.mutex.Lock()
	c.entries[key] = Entry{
		expires: time.Now().Add(c.entryExpirationDefault),
		content: value,
	}
	c.mutex.Unlock()
}

// Get returns the value (if found) and a bool
// indicating whether the value exists
func (c Cache) Get(key string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	entry, exists := c.entries[key]
	if !exists || entry.expired() {
		return "", false
	}

	return entry.content, true
}
