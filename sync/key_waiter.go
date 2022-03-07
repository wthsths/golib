package sync

import (
	go_sync "sync"
)

// KeyWaiter is used for blocking simultaneous usage of keys referring to same entities.
type KeyWaiter struct {
	mu       go_sync.RWMutex
	waitKeys map[string]*waitKey
}

// NewKeyWaiter creates a new waiter instance.
// Instance must be shared across application wherever a shared mutex is necessary for the same named entity set.
func NewKeyWaiter() *KeyWaiter {
	return &KeyWaiter{waitKeys: make(map[string]*waitKey)}
}

// waitKey wraps a single mutex with the waiting count to be used for each unique key.
type waitKey struct {
	go_sync.Mutex
	waiting int
}

// Acquire locks the given key for current key waiter across the application.
func (g *KeyWaiter) Acquire(key string) {
	g.mu.Lock() // Lock to ensure keys are not read before current change is made
	if g.waitKeys == nil {
		g.waitKeys = make(map[string]*waitKey)
	}

	keyGroup, ok := g.waitKeys[key]
	if ok {
		// Get into the waiting call count for the key
		keyGroup.waiting++
		g.mu.Unlock() // Unlock while waiting so others can acquire

		// Wait for mutex to be released and acquire
		keyGroup.Lock()

		g.mu.Lock()
		// Get off the waiting call count
		keyGroup.waiting--
	} else {
		// If key is not in use, create a new waitKey
		keyGroup = new(waitKey)
		g.waitKeys[key] = keyGroup
		// Acquire key
		keyGroup.Lock()
	}

	g.mu.Unlock()
}

// Release unlocks the given key for current key waiter across the application.
// Panics on releasing an already released or non-exiting key.
func (g *KeyWaiter) Release(key string) {
	g.mu.Lock()

	// Make sure key exist
	keyGroup, ok := g.waitKeys[key]
	if ok {
		// Release mutex for the key
		keyGroup.Unlock()
	} else {
		g.mu.Unlock()
		panic("KeyWaiter: tried to release a non existing key")
	}

	// If no call is waiting to acquire, delete key
	if keyGroup.waiting <= 0 {
		// glog.Infof("deleting key %s", key)
		delete(g.waitKeys, key)
	}

	g.mu.Unlock()
}
