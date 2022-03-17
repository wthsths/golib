package sync

import "sync"

// SyncedBool is goroutine -afe variant of boolean type.
type SyncedBool struct {
	value bool
	mutex *sync.Mutex
}

// NewSyncedBool creates a new goroutine-safe value of boolean type.
func NewSyncedBool(initialValue bool) *SyncedBool {
	return &SyncedBool{
		value: initialValue,
		mutex: &sync.Mutex{},
	}
}

// Get returns wrapped boolean value in a goroutine-safe manner.
func (sb *SyncedBool) Get() bool {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	return sb.value
}

// Set updates the wrapped boolean value in a goroutine-safe manner.
func (sb *SyncedBool) Set(value bool) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	sb.value = value
}
