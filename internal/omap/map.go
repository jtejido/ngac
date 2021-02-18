package omap

type OrderedMap interface {
	// Add adds items to the map or replace an existing key.
	Add(key interface{}, value interface{})

	// Add items from another map to this instance. Updates value if existing.
	AddMap(m OrderedMap)

	// Get returns the value of a key from the map.
	Get(key interface{}) (value interface{}, found bool)

	// Remove deletes a key-value pair from the map.
	Remove(key interface{})

	// Len return the map number of key-value pairs.
	Len() int

	// Keys return the keys in the map in insertion order.
	Keys() []interface{}

	// Values return the values in the map in insertion order.
	Values() []interface{}

	Clone() OrderedMap
}

func NewOrderedMap() OrderedMap {
	m := newSafeOrderedMap()
	return &m
}

// NewThreadUnsafeSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewUnsafeOrderedMap() OrderedMap {
	return newUnsafeOrderedMap()
}
