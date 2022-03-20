package gl_maps

// OrderedStrMap wraps a map and a slice to keep track of insertion order of the map entries.
//
// Loop through map with OrderedKeys() for insertion-time-ordered access.
type OrderedStrMap struct {
	innerMap map[string]string
	keys     []string
}

func NewOrderedStrMap(initialCap int) *OrderedStrMap {
	return &OrderedStrMap{
		innerMap: make(map[string]string, initialCap),
		keys:     make([]string, 0, initialCap),
	}
}

func (m *OrderedStrMap) Get(key string) string {
	return m.innerMap[key]
}

// GetWithCheck will return a second value which indicates whether the key exists or not.
func (m *OrderedStrMap) GetWithCheck(key string) (string, bool) {
	value, ok := m.innerMap[key]
	return value, ok
}

func (m *OrderedStrMap) Set(key, value string) {
	m.keys = append(m.keys, key)
	m.innerMap[key] = value
}

// OrderedKeys returns insertion-time-ordered COPY of the map keys.
func (m *OrderedStrMap) OrderedKeys() []string {
	return m.keys
}

// CopyInnerMap returns a COPY of the inner map.
func (m *OrderedStrMap) CopyInnerMap() map[string]string {
	clone := make(map[string]string, len(m.innerMap))
	for k, v := range m.innerMap {
		clone[k] = v
	}
	return clone
}
