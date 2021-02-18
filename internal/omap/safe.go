package omap

import (
	"sync"
)

var (
	sm *safeOrderedMap
	_  OrderedMap = sm
)

type safeOrderedMap struct {
	m *unsafeOrderedMap
	sync.RWMutex
}

func newSafeOrderedMap() safeOrderedMap {
	return safeOrderedMap{m: newUnsafeOrderedMap()}
}

func (m *safeOrderedMap) Add(key interface{}, value interface{}) {
	m.Lock()
	m.m.Add(key, value)
	m.Unlock()
}

func (m *safeOrderedMap) AddMap(om OrderedMap) {
	o := om.(*safeOrderedMap)
	m.Lock()
	o.Lock()
	m.m.AddMap(o.m)
	m.Unlock()
	o.Unlock()
}

func (m *safeOrderedMap) Get(key interface{}) (value interface{}, found bool) {
	m.RLock()
	value, found = m.m.Get(key)
	m.RUnlock()
	return
}

func (m *safeOrderedMap) Remove(key interface{}) {
	m.Lock()
	defer m.Unlock()

	// Check key exists
	if _, found := m.m.store[key]; !found {
		return
	}

	// Remove the value from the store
	delete(m.m.store, key)

	// Remove the key
	for i := range m.m.keys {
		if m.m.keys[i] == key {
			m.m.keys = append(m.m.keys[:i], m.m.keys[i+1:]...)
			break
		}
	}
}

func (m *safeOrderedMap) Len() int {
	m.Lock()
	defer m.Unlock()

	return len(m.m.store)
}

func (m *safeOrderedMap) Keys() []interface{} {
	m.RLock()
	defer m.RUnlock()

	return m.m.keys
}

func (m *safeOrderedMap) Values() []interface{} {
	m.RLock()
	defer m.RUnlock()

	values := make([]interface{}, len(m.m.store))

	for i, key := range m.m.keys {
		values[i] = m.m.store[key]
	}

	return values
}

func (m *safeOrderedMap) Clone() OrderedMap {
	m.RLock()

	unsafeClone := m.m.Clone().(*unsafeOrderedMap)
	ret := &safeOrderedMap{m: unsafeClone}
	m.RUnlock()

	return ret
}
