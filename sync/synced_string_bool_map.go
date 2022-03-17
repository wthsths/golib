package sync

import go_sync "sync"

// SyncedStringBoolMap is goroutine-safe variant of map[string]bool type.
type SyncedStringBoolMap struct {
	mutex    go_sync.Mutex
	innerMap map[string]bool
}

// NewSyncedStringBoolMap creates a new goroutine-safe value of map[string]bool type.
func NewSyncedStringBoolMap(initialCap int) *SyncedStringBoolMap {
	return &SyncedStringBoolMap{
		innerMap: make(map[string]bool, initialCap),
	}
}

// Get returns wrapped map[string]bool value in a goroutine-safe manner.
func (m *SyncedStringBoolMap) Get(key string) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.innerMap[key]
}

func (m *SyncedStringBoolMap) GetWithCheck(key string) (bool, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	value, ok := m.innerMap[key]
	return value, ok
}

func (m *SyncedStringBoolMap) GetKeys() []string {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	keys := make([]string, 0, len(m.innerMap))

	for k := range m.innerMap {
		keys = append(keys, k)
	}
	return keys
}

// Set updates the wrapped map[string]bool value in a goroutine-safe manner.
func (m *SyncedStringBoolMap) Set(key string, value bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.innerMap[key] = value
}
