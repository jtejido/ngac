package set

import "sync"

var (
	ss *safeSet
	_  Set = ss
)

type safeSet struct {
	s *unsafeSet
	sync.RWMutex
}

func newSafeSet() safeSet {
	return safeSet{s: newUnsafeSet()}
}

func (s *safeSet) Add(i ...interface{}) {
	s.Lock()
	s.s.Add(i...)
	s.Unlock()
}

func (s *safeSet) Contains(i ...interface{}) bool {
	s.RLock()
	ret := s.s.Contains(i...)
	s.RUnlock()
	return ret
}

func (s *safeSet) IsSubset(other Set) bool {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	ret := s.s.IsSubset(o.s)
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet) IsProperSubset(other Set) bool {
	o := other.(*safeSet)

	s.RLock()
	defer s.RUnlock()
	o.RLock()
	defer o.RUnlock()

	return s.s.IsProperSubset(o.s)
}

func (s *safeSet) IsSuperset(other Set) bool {
	return other.IsSubset(s)
}

func (s *safeSet) IsProperSuperset(other Set) bool {
	return other.IsProperSubset(s)
}

func (s *safeSet) Union(other Set) Set {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	ret := &safeSet{s: s.s.Union(o.s).(*unsafeSet)}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet) AddFrom(other Set) {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	s.s.AddFrom(o.s)

}

func (s *safeSet) Intersect(other Set) Set {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	ret := &safeSet{s: s.s.Intersect(o.s).(*unsafeSet)}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet) RetainFrom(other Set) {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	s.s.RetainFrom(o.s)
}

func (s *safeSet) IsDisjointFrom(other Set) bool {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	intersectSet := s.s.Intersect(o.s).(*unsafeSet)
	isDisjoint := intersectSet.Len() == 0

	s.RUnlock()
	o.RUnlock()

	return isDisjoint
}

func (s *safeSet) Difference(other Set) Set {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	ret := &safeSet{s: s.s.Difference(o.s).(*unsafeSet)}
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet) RemoveFrom(other Set) {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()
	defer s.RUnlock()
	defer o.RUnlock()

	s.s.RemoveFrom(o.s)
}

func (s *safeSet) Clear() {
	s.Lock()
	s.s = newUnsafeSet()
	s.Unlock()
}

func (s *safeSet) Remove(i interface{}) {
	s.Lock()
	defer s.Unlock()
	s.s.Remove(i)

}

func (s *safeSet) Len() int {
	s.RLock()
	defer s.RUnlock()
	return len(s.s.m)
}

func (s *safeSet) Iter() <-chan interface{} {
	ch := make(chan interface{})
	go func() {
		s.RLock()
		for _, key := range s.s.keys {
			ch <- s.s.m[key]
		}
		close(ch)
		s.RUnlock()
	}()

	return ch
}

func (s *safeSet) Iterator() SetIterator {
	return newSafeIterator(s)
}

func (s *safeSet) Equal(other Set) bool {
	o := other.(*safeSet)

	s.RLock()
	o.RLock()

	ret := s.s.Equal(o.s)
	s.RUnlock()
	o.RUnlock()
	return ret
}

func (s *safeSet) Clone() Set {
	s.RLock()

	unsafeClone := s.s.Clone().(*unsafeSet)
	ret := &safeSet{s: unsafeClone}
	s.RUnlock()
	return ret
}

func (s *safeSet) ToSlice() []interface{} {
	keys := make([]interface{}, 0, s.Len())
	s.RLock()
	for _, key := range s.s.keys {
		keys = append(keys, s.s.m[key])
	}
	s.RUnlock()
	return keys
}

func (s *safeSet) Filter(filter func(interface{}) bool) {
	s.RLock()
	s.s.Filter(filter)
	s.RUnlock()
}
