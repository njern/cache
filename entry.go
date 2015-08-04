package cache

import "time"

// An Entry encapsulates a value stored in a Cache
type Entry struct {
	expires time.Time
	content string
}

// expired returns true if the entry has expired
func (e Entry) expired() bool {
	return time.Now().After(e.expires)
}
