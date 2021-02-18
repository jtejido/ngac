package set

// ------ Stateful iterator ------
type SetIterator interface {
	HasNext() bool     // checks for available rows
	Next() interface{} // Advances to the next row in the set
}

type safeIterator struct {
	*safeSet
	index int
}

func newSafeIterator(s *safeSet) *safeIterator {
	return &safeIterator{s, -1}
}

// Next returns the next item in the collection.
func (i *safeIterator) Next() interface{} {
	i.RLock()
	defer i.RUnlock()

	i.index++
	item := i.s.m[i.s.keys[i.index]]

	return item
}

// HasNext return true if there are values to be read.
func (i *safeIterator) HasNext() bool {
	i.RLock()
	defer i.RUnlock()
	return i.index < (len(i.s.m) - 1)
}

type unsafeIterator struct {
	*unsafeSet
	index int
}

func newUnsafeIterator(s *unsafeSet) *unsafeIterator {
	return &unsafeIterator{s, -1}
}

// Next returns the next item in the collection.
func (i *unsafeIterator) Next() interface{} {
	i.index++
	item := i.m[i.keys[i.index]]
	return item
}

// HasNext return true if there are values to be read.
func (i *unsafeIterator) HasNext() bool {
	return i.index < (len(i.m) - 1)
}
