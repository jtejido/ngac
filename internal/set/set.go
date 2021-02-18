package set

// Ordered Set Interface
// All values are iterated in the order they're added in.
//
// Take note that all operations work within the confine
// of a shared type used between the items of the sets.
// Thus, a 'string' Set shall not be compared with an
// 'int64' set for instance, or it will cause panic.
type Set interface {
	Add(i ...interface{})
	Len() int
	Clear()

	/*
	 *
	 * Below methods are equality checks.
	 *
	 */

	// Checks if this set contains the value specified
	Contains(i ...interface{}) bool

	// Determines if two sets are equal to each
	// other. If they have the same cardinality
	// and contain the same elements, they are
	// considered equal. The order in which
	// the elements were added is irrelevant.
	Equal(other Set) bool

	// Determines if this set, which intersects with the
	// other set, is an empty set.
	IsDisjointFrom(other Set) bool

	// Determines if every element in this set is in
	// the other set but the two sets are not equal.
	IsProperSubset(other Set) bool

	// Determines if every element in the other set
	// is in this set but the two sets are not
	// equal.
	IsProperSuperset(other Set) bool

	// Determines if every element in this set is in
	// the other set.
	IsSubset(other Set) bool

	// Determines if every element in the other set
	// is in this set.
	IsSuperset(other Set) bool

	// Iterates over all item in the order they're
	// added.
	Iter() <-chan interface{}

	Iterator() SetIterator
	Remove(i interface{})
	ToSlice() []interface{}

	/*
	 *
	 * Below methods produces a new Set instance. (see in-place operations. e.g., AddFrom, RetainFrom, and RemoveFrom)
	 *
	 */

	// Returns a new set with all elements in both sets.
	Union(other Set) Set

	// Returns a new set containing only the elements
	// that exist only in both sets.
	Intersect(other Set) Set

	// Returns the difference between this set
	// and other. The returned set will contain
	// all elements of this set that are not in
	// the other set.
	Difference(other Set) Set

	// Clones the values from this set to a new set instance.
	Clone() Set

	/*
	 *
	 * Below methods are in-place methods.
	 *
	 */

	// Add items from the other set to this set (if non-existent to this set).
	// This operation is same as Add(i ...interface{}) but meant to be easy
	// for other set implementing same interface.
	AddFrom(other Set)

	// Retains the elements that exist only in both sets.
	RetainFrom(other Set)

	// It removes all elements in this set that is/are present in the other set.
	RemoveFrom(other Set)

	Filter(func(interface{}) bool)
}

type Equaler interface {
	Equals(interface{}) bool
}

func NewSet(s ...interface{}) Set {
	set := newSafeSet()
	for _, item := range s {
		set.Add(item)
	}
	return &set
}

// NewSetFromSlice creates and returns a reference to a set from an
// existing slice.  Operations on the resulting set are thread-safe.
func NewSetFromSlice(s []interface{}) Set {
	a := NewSet(s...)
	return a
}

// NewThreadUnsafeSet creates and returns a reference to an empty set.
// Operations on the resulting set are not thread-safe.
func NewUnsafeSet(s ...interface{}) Set {
	set := newUnsafeSet()
	for _, item := range s {
		set.Add(item)
	}
	return set
}

// NewThreadUnsafeSetFromSlice creates and returns a reference to a
// set from an existing slice.  Operations on the resulting set are
// not thread-safe.
func NewUnsafeSetFromSlice(s []interface{}) Set {
	a := newUnsafeSet()
	for _, item := range s {
		a.Add(item)
	}
	return a
}
