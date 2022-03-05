package sync

import go_sync "sync"

type SyncedInt64BoolMap struct {
	mutex    go_sync.Mutex
	innerMap map[int64]bool
}

func NewSyncedInt64BoolMap(initialCap int) *SyncedInt64BoolMap {
	return &SyncedInt64BoolMap{
		innerMap: make(map[int64]bool, initialCap),
	}
}

func (m *SyncedInt64BoolMap) Set(key int64, value bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.innerMap[key] = value
}

func (m *SyncedInt64BoolMap) Get(key int64) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.innerMap[key]
}

func (m *SyncedInt64BoolMap) GetWithCheck(key int64) (bool, bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	value, ok := m.innerMap[key]
	return value, ok
}

func (m *SyncedInt64BoolMap) GetKeys() []int64 {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	keys := make([]int64, 0, len(m.innerMap))

	for k := range m.innerMap {
		keys = append(keys, k)
	}
	return keys
}
