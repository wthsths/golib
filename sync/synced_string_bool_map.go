package sync

import go_sync "sync"

type SyncedStringBoolMap struct {
	mutex    go_sync.Mutex
	innerMap map[string]bool
}

func NewSyncedStringBoolMap(initialCap int) *SyncedStringBoolMap {
	return &SyncedStringBoolMap{
		innerMap: make(map[string]bool, initialCap),
	}
}

func (m *SyncedStringBoolMap) Set(key string, value bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.innerMap[key] = value
}

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
