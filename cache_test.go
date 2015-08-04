package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(nil, nil)
	key := "1"
	value := "1"
	otherKey := "2"

	c.Set(key, value)

	// Check that the real key exists and returns the right result
	entry, exists := c.Get(key)
	if entry != value {
		t.Error("a Get() should return the correct value for the key")
	}
	if exists == false {
		t.Error("exists should be true if the item exists")
	}

	// Check that the non-existent key returns the right result
	entry, exists = c.Get(otherKey)
	if entry != "" {
		t.Error("a non-existent key should return the empty string")
	}
	if exists == true {
		t.Error("exists should be false if the item does not exist")
	}
}

func TestCacheEviction(t *testing.T) {
	shortExpirationTime := time.Millisecond
	c := New(&shortExpirationTime, &shortExpirationTime)
	key := "1"
	value := "1"
	c.Set(key, value)

	// Check that the key is evicted once timed out
	time.Sleep(time.Millisecond)
	entry, exists := c.Get(key)
	if entry != "" {
		t.Error("a non-existent key should return the empty string")
	}
	if exists == true {
		t.Error("exists should be false if the item does not exist")
	}
}
