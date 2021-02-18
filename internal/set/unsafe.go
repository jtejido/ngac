package set

var (
	us *unsafeSet
	_  Set = us
)

type unsafeSet struct {
	currentIndex int
	m            map[int]interface{}
	index        map[interface{}]int
	keys         []int
}

func newUnsafeSet() *unsafeSet {
	return &unsafeSet{currentIndex: 0, m: make(map[int]interface{}), keys: make([]int, 0), index: make(map[interface{}]int)}
}

func (s *unsafeSet) Add(i ...interface{}) {
	for _, item := range i {
		var found bool
		// check if equaler, if not proceed with interface checking
		if e1, ok := item.(Equaler); ok {
			for k, _ := range s.index {
				if e2, ok2 := k.(Equaler); ok2 {
					if e1.Equals(e2) {
						found = true
						break
					}
				}
			}
		} else {
			_, found = s.index[item]
		}

		if found {
			continue
		}

		if _, ok := s.m[s.currentIndex]; !ok {
			s.keys = append(s.keys, s.currentIndex)
		}

		s.m[s.currentIndex] = item
		s.index[item] = s.currentIndex
		s.currentIndex++
	}
}

func (s *unsafeSet) Contains(i ...interface{}) bool {
	for _, val := range i {
		var found bool
		// check if equaler, if not proceed with interface checking
		if e1, ok := val.(Equaler); ok {
			for k, _ := range s.index {
				if e2, ok2 := k.(Equaler); ok2 {
					if e1.Equals(e2) {
						found = true
						break
					}
				}
			}
		} else {
			_, found = s.index[val]
		}

		if !found {
			return false
		}
	}
	return true
}

func (s *unsafeSet) IsSubset(other Set) bool {
	_ = other.(*unsafeSet)
	if s.Len() > other.Len() {
		return false
	}
	for elem := range s.index {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (s *unsafeSet) IsProperSubset(other Set) bool {
	return s.IsSubset(other) && !s.Equal(other)
}

func (s *unsafeSet) IsSuperset(other Set) bool {
	return other.IsSubset(s)
}

func (s *unsafeSet) IsProperSuperset(other Set) bool {
	return s.IsSuperset(other) && !s.Equal(other)
}

func (s *unsafeSet) IsDisjointFrom(other Set) bool {
	o := other.(*unsafeSet)
	intersectSet := s.Intersect(o).(*unsafeSet)

	return intersectSet.Len() == 0
}

func (s *unsafeSet) Clear() {
	*s = *newUnsafeSet()
}

func (s *unsafeSet) Remove(i interface{}) {
	var found bool
	var index int
	// check if equaler, if not proceed with interface checking
	if e1, ok := i.(Equaler); ok {
		var k interface{}
		for k, index = range s.index {
			if e2, ok2 := k.(Equaler); ok2 {
				if e1.Equals(e2) {
					found = true
					break
				}
			}
		}
	} else {
		index, found = s.index[i]
	}

	if !found {
		return
	}

	// Check key exists
	if _, found := s.m[index]; !found {
		return
	}

	delete(s.m, index)

	// Remove the key
	for i := range s.keys {
		if s.keys[i] == index {
			s.keys = append(s.keys[:i], s.keys[i+1:]...)
			break
		}
	}

	delete(s.index, i)
}

func (s *unsafeSet) Len() int {
	return len(s.m)
}

func (s *unsafeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		for _, key := range s.keys {
			ch <- s.m[key]
		}

		close(ch)
	}()

	return ch
}

func (s *unsafeSet) Iterator() SetIterator {
	return newUnsafeIterator(s)
}

func (s *unsafeSet) Equal(other Set) bool {
	_ = other.(*unsafeSet)

	if s.Len() != other.Len() {
		return false
	}
	for elem := range s.index {
		if !other.Contains(elem) {
			return false
		}
	}
	return true
}

func (s *unsafeSet) Clone() Set {
	clonedSet := newUnsafeSet()
	for _, key := range s.keys {
		clonedSet.Add(s.m[key])
	}
	return clonedSet
}

func (s *unsafeSet) ToSlice() []interface{} {
	keys := make([]interface{}, 0, s.Len())
	for _, key := range s.keys {
		keys = append(keys, s.m[key])
	}

	return keys
}

func (s *unsafeSet) Union(other Set) Set {
	o := other.(*unsafeSet)

	unionedSet := newUnsafeSet()

	for elem := range s.index {
		unionedSet.Add(elem)
	}

	for elem := range o.index {
		unionedSet.Add(elem)
	}

	return unionedSet
}

func (s *unsafeSet) AddFrom(other Set) {
	o := other.(*unsafeSet)

	for elem := range o.index {
		if !s.Contains(elem) {
			s.Add(elem)
		}
	}
}

func (s *unsafeSet) Intersect(other Set) Set {
	o := other.(*unsafeSet)

	intersection := newUnsafeSet()
	// loop over smaller set
	if s.Len() < other.Len() {
		for elem := range s.index {
			if other.Contains(elem) {
				intersection.Add(elem)
			}
		}
	} else {
		for elem := range o.index {
			if s.Contains(elem) {
				intersection.Add(elem)
			}
		}
	}
	return intersection
}

func (s *unsafeSet) RetainFrom(other Set) {

	for elem := range s.index {
		if !other.Contains(elem) {
			s.Remove(elem)
		}
	}
}

func (s *unsafeSet) Difference(other Set) Set {
	_ = other.(*unsafeSet)

	difference := newUnsafeSet()
	for elem := range s.index {
		if !other.Contains(elem) {
			difference.Add(elem)
		}
	}

	return difference
}

func (s *unsafeSet) RemoveFrom(other Set) {
	for elem := range s.index {
		if other.Contains(elem) {
			s.Remove(elem)
		}
	}
}

func (s *unsafeSet) Filter(filter func(interface{}) bool) {
	for elem := range s.index {
		if filter(elem) {
			s.Remove(elem)
		}
	}
}
