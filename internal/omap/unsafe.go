package omap

var (
	um *unsafeOrderedMap
	_  OrderedMap = um
)

type unsafeOrderedMap struct {
	keys  []interface{}
	store map[interface{}]interface{}
}

func newUnsafeOrderedMap() *unsafeOrderedMap {
	m := &unsafeOrderedMap{
		keys:  make([]interface{}, 0),
		store: make(map[interface{}]interface{}),
	}

	return m
}

func (m *unsafeOrderedMap) Add(key interface{}, value interface{}) {
	if _, ok := m.store[key]; !ok {
		m.keys = append(m.keys, key)
	}

	m.store[key] = value
}

func (m *unsafeOrderedMap) AddMap(om OrderedMap) {
	o := om.(*unsafeOrderedMap)

	for _, key := range o.keys {
		otherVal, _ := o.store[key]
		if _, ok := m.store[key]; !ok {
			m.keys = append(m.keys, key)
		}

		m.store[key] = otherVal
	}
}

func (m *unsafeOrderedMap) Get(key interface{}) (value interface{}, found bool) {
	value, found = m.store[key]
	return value, found
}

func (m *unsafeOrderedMap) Remove(key interface{}) {
	if _, found := m.store[key]; !found {
		return
	}

	delete(m.store, key)

	for i := range m.keys {
		if m.keys[i] == key {
			m.keys = append(m.keys[:i], m.keys[i+1:]...)
			break
		}
	}
}

func (m *unsafeOrderedMap) Len() int {
	return len(m.store)
}

func (m *unsafeOrderedMap) Keys() []interface{} {
	return m.keys
}

func (m *unsafeOrderedMap) Values() []interface{} {
	values := make([]interface{}, len(m.store))

	for i, key := range m.keys {
		values[i] = m.store[key]
	}

	return values
}

func (m *unsafeOrderedMap) Clone() OrderedMap {
	clonedSet := newUnsafeOrderedMap()
	for _, key := range m.keys {
		clonedSet.Add(key, m.store[key])
	}
	return clonedSet
}
