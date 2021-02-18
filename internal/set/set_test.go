package set

import "testing"

func makeSet(ints []int) Set {
	set := NewSet()
	for _, i := range ints {
		set.Add(i)
	}
	return set
}

func makeUnsafeSet(ints []int) Set {
	set := NewUnsafeSet()
	for _, i := range ints {
		set.Add(i)
	}
	return set
}

func assertEqual(a, b Set, t *testing.T) {
	if !a.Equal(b) {
		t.Errorf("%v != %v\n", a, b)
	}
}

func TestNewSet(t *testing.T) {
	a := NewSet()
	if a.Len() != 0 {
		t.Error("NewSet should start out as an empty set")
	}

	assertEqual(NewSetFromSlice([]interface{}{}), NewSet(), t)
	assertEqual(NewSetFromSlice([]interface{}{1}), NewSet(1), t)
	assertEqual(NewSetFromSlice([]interface{}{1, 2}), NewSet(2, 1), t)
	assertEqual(NewSetFromSlice([]interface{}{"a"}), NewSet("a"), t)
	assertEqual(NewSetFromSlice([]interface{}{"a", "b"}), NewSet("b", "a"), t)
}

func TestNewUnsafeSet(t *testing.T) {
	a := NewUnsafeSet()

	if a.Len() != 0 {
		t.Error("NewSet should start out as an empty set")
	}
}

func TestAddSet(t *testing.T) {
	a := makeSet([]int{1, 2, 3})

	if a.Len() != 3 {
		t.Error("AddSet does not have a size of 3 even though 3 items were added to a new set")
	}
}

func TestAddUnsafeSet(t *testing.T) {
	a := makeUnsafeSet([]int{1, 2, 3})

	if a.Len() != 3 {
		t.Error("AddSet does not have a size of 3 even though 3 items were added to a new set")
	}
}

func TestAddSetNoDuplicate(t *testing.T) {
	a := makeSet([]int{7, 5, 3, 7})

	if a.Len() != 3 {
		t.Error("AddSetNoDuplicate set should have 3 elements since 7 is a duplicate")
	}

	if !(a.Contains(7) && a.Contains(5) && a.Contains(3)) {
		t.Error("AddSetNoDuplicate set should have a 7, 5, and 3 in it.")
	}
}

func TestAddUnsafeSetNoDuplicate(t *testing.T) {
	a := makeUnsafeSet([]int{7, 5, 3, 7})

	if a.Len() != 3 {
		t.Error("AddSetNoDuplicate set should have 3 elements since 7 is a duplicate")
	}

	if !(a.Contains(7) && a.Contains(5) && a.Contains(3)) {
		t.Error("AddSetNoDuplicate set should have a 7, 5, and 3 in it.")
	}
}

func TestRemoveSet(t *testing.T) {
	a := makeSet([]int{6, 3, 1})

	a.Remove(3)

	if a.Len() != 2 {
		t.Error("RemoveSet should only have 2 items in the set")
	}

	if !(a.Contains(6) && a.Contains(1)) {
		t.Error("RemoveSet should have only items 6 and 1 in the set")
	}

	a.Remove(6)
	a.Remove(1)

	if a.Len() != 0 {
		t.Error("RemoveSet should be an empty set after removing 6 and 1")
	}
}

func TestRemoveUnsafeSet(t *testing.T) {
	a := makeUnsafeSet([]int{6, 3, 1})

	a.Remove(3)

	if a.Len() != 2 {
		t.Error("RemoveSet should only have 2 items in the set")
	}

	if !(a.Contains(6) && a.Contains(1)) {
		t.Error("RemoveSet should have only items 6 and 1 in the set")
	}

	a.Remove(6)
	a.Remove(1)

	if a.Len() != 0 {
		t.Error("RemoveSet should be an empty set after removing 6 and 1")
	}
}

func TestContainsSet(t *testing.T) {
	a := NewSet()

	a.Add(71)

	if !a.Contains(71) {
		t.Error("ContainsSet should contain 71")
	}

	a.Remove(71)

	if a.Contains(71) {
		t.Error("ContainsSet should not contain 71")
	}

	a.Add(13)
	a.Add(7)
	a.Add(1)

	if !(a.Contains(13) && a.Contains(7) && a.Contains(1)) {
		t.Error("ContainsSet should contain 13, 7, 1")
	}
}

func TestContainsUnsafeSet(t *testing.T) {
	a := NewUnsafeSet()

	a.Add(71)

	if !a.Contains(71) {
		t.Error("ContainsSet should contain 71")
	}

	a.Remove(71)

	if a.Contains(71) {
		t.Error("ContainsSet should not contain 71")
	}

	a.Add(13)
	a.Add(7)
	a.Add(1)

	if !(a.Contains(13) && a.Contains(7) && a.Contains(1)) {
		t.Error("ContainsSet should contain 13, 7, 1")
	}
}

func TestContainsMultipleSet(t *testing.T) {
	a := makeSet([]int{8, 6, 7, 5, 3, 0, 9})

	if !a.Contains(8, 6, 7, 5, 3, 0, 9) {
		t.Error("ContainsAll should contain Jenny's phone number")
	}

	if a.Contains(8, 6, 11, 5, 3, 0, 9) {
		t.Error("ContainsAll should not have all of these numbers")
	}
}

func TestContainsMultipleUnsafeSet(t *testing.T) {
	a := makeUnsafeSet([]int{8, 6, 7, 5, 3, 0, 9})

	if !a.Contains(8, 6, 7, 5, 3, 0, 9) {
		t.Error("ContainsAll should contain Jenny's phone number")
	}

	if a.Contains(8, 6, 11, 5, 3, 0, 9) {
		t.Error("ContainsAll should not have all of these numbers")
	}
}

func TestClearSet(t *testing.T) {
	a := makeSet([]int{2, 5, 9, 10})

	a.Clear()

	if a.Len() != 0 {
		t.Error("ClearSet should be an empty set")
	}
}

func TestClearUnsafeSet(t *testing.T) {
	a := makeUnsafeSet([]int{2, 5, 9, 10})

	a.Clear()

	if a.Len() != 0 {
		t.Error("ClearSet should be an empty set")
	}
}

func TestLenSet(t *testing.T) {
	a := NewSet()

	if a.Len() != 0 {
		t.Error("set should be an empty set")
	}

	a.Add(1)

	if a.Len() != 1 {
		t.Error("set should have a size of 1")
	}

	a.Remove(1)

	if a.Len() != 0 {
		t.Error("set should be an empty set")
	}

	a.Add(9)

	if a.Len() != 1 {
		t.Error("set should have a size of 1")
	}

	a.Clear()

	if a.Len() != 0 {
		t.Error("set should have a size of 1")
	}
}

func TestLenUnsafeSet(t *testing.T) {
	a := NewUnsafeSet()

	if a.Len() != 0 {
		t.Error("set should be an empty set")
	}

	a.Add(1)

	if a.Len() != 1 {
		t.Error("set should have a size of 1")
	}

	a.Remove(1)

	if a.Len() != 0 {
		t.Error("set should be an empty set")
	}

	a.Add(9)

	if a.Len() != 1 {
		t.Error("set should have a size of 1")
	}

	a.Clear()

	if a.Len() != 0 {
		t.Error("set should have a size of 1")
	}
}

func TestSetIsSubset(t *testing.T) {
	a := makeSet([]int{1, 2, 3, 5, 7})

	b := NewSet()
	b.Add(3)
	b.Add(5)
	b.Add(7)

	if !b.IsSubset(a) {
		t.Error("set b should be a subset of set a")
	}

	b.Add(72)

	if b.IsSubset(a) {
		t.Error("set b should not be a subset of set a because it contains 72 which is not in the set of a")
	}
}

func TestSetIsProperSubset(t *testing.T) {
	a := makeSet([]int{1, 2, 3, 5, 7})
	b := makeSet([]int{7, 5, 3, 2, 1})

	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	if a.IsProperSubset(b) {
		t.Error("set a should not be a proper subset of set b (they're equal)")
	}

	b.Add(72)

	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	if !a.IsProperSubset(b) {
		t.Error("set a should be a proper subset of set b")
	}
}

func TestUnsafeSetIsSubset(t *testing.T) {
	a := makeUnsafeSet([]int{1, 2, 3, 5, 7})

	b := NewUnsafeSet()
	b.Add(3)
	b.Add(5)
	b.Add(7)

	if !b.IsSubset(a) {
		t.Error("set b should be a subset of set a")
	}

	b.Add(72)

	if b.IsSubset(a) {
		t.Error("set b should not be a subset of set a because it contains 72 which is not in the set of a")
	}
}

func TestUnsafeSetIsProperSubset(t *testing.T) {
	a := makeUnsafeSet([]int{1, 2, 3, 5, 7})
	b := NewUnsafeSet()
	b.Add(7)
	b.Add(1)
	b.Add(5)
	b.Add(3)
	b.Add(2)

	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	if a.IsProperSubset(b) {
		t.Error("set a should not be a proper subset of set b (they're equal)")
	}

	b.Add(72)

	if !a.IsSubset(b) {
		t.Error("set a should be a subset of set b")
	}
	if !a.IsProperSubset(b) {
		t.Error("set a should be a proper subset of set b because set b has 72")
	}
}

func TestSetIsSuperset(t *testing.T) {
	a := NewSet()
	a.Add(9)
	a.Add(5)
	a.Add(2)
	a.Add(1)
	a.Add(11)

	b := NewSet()
	b.Add(5)
	b.Add(2)
	b.Add(11)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set b has a 42")
	}
}

func TestSetIsProperSuperset(t *testing.T) {
	a := NewSet()
	a.Add(5)
	a.Add(2)
	a.Add(11)

	b := NewSet()
	b.Add(2)
	b.Add(5)
	b.Add(11)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b (they're equal)")
	}

	a.Add(9)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if !a.IsProperSuperset(b) {
		t.Error("set a not be a proper superset of set b because set a has a 9")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set b has a 42")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b because set b has a 42")
	}
}

func TestUnsafeSetIsSuperset(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(9)
	a.Add(5)
	a.Add(2)
	a.Add(1)
	a.Add(11)

	b := NewUnsafeSet()
	b.Add(5)
	b.Add(2)
	b.Add(11)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set a has a 42")
	}
}

func TestUnsafeSetIsProperSuperset(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(5)
	a.Add(2)
	a.Add(11)

	b := NewUnsafeSet()
	b.Add(2)
	b.Add(5)
	b.Add(11)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b (they're equal)")
	}

	a.Add(9)

	if !a.IsSuperset(b) {
		t.Error("set a should be a superset of set b")
	}
	if !a.IsProperSuperset(b) {
		t.Error("set a not be a proper superset of set b because set a has a 9")
	}

	b.Add(42)

	if a.IsSuperset(b) {
		t.Error("set a should not be a superset of set b because set b has a 42")
	}
	if a.IsProperSuperset(b) {
		t.Error("set a should not be a proper superset of set b because set b has a 42")
	}
}

func TestSetUnion(t *testing.T) {
	a := NewSet()

	b := NewSet()
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	b.Add(5)

	c := a.Union(b)

	if c.Len() != 5 {
		t.Error("set c is unioned with an empty set and therefore should have 5 elements in it")
	}

	d := NewSet()
	d.Add(10)
	d.Add(14)
	d.Add(0)

	e := c.Union(d)
	if e.Len() != 8 {
		t.Error("set e should should have 8 elements in it after being unioned with set c to d")
	}

	f := NewSet()
	f.Add(14)
	f.Add(3)

	g := f.Union(e)
	if g.Len() != 8 {
		t.Error("set g should still have 8 elements in it after being unioned with set f that has duplicates")
	}
}

func TestUnsafeSetUnion(t *testing.T) {
	a := NewUnsafeSet()

	b := NewUnsafeSet()
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	b.Add(5)

	c := a.Union(b)

	if c.Len() != 5 {
		t.Error("set c is unioned with an empty set and therefore should have 5 elements in it")
	}

	d := NewUnsafeSet()
	d.Add(10)
	d.Add(14)
	d.Add(0)

	e := c.Union(d)
	if e.Len() != 8 {
		t.Error("set e should should have 8 elements in it after being unioned with set c to d")
	}

	f := NewUnsafeSet()
	f.Add(14)
	f.Add(3)

	g := f.Union(e)
	if g.Len() != 8 {
		t.Error("set g should still have 8 elements in it after being unioned with set f that has duplicates")
	}
}

func TestUnsafeSetAddAsSlice(t *testing.T) {
	a := NewUnsafeSet()
	b := NewUnsafeSet()
	b.Add(1)
	b.Add(2)
	b.Add(3)
	b.Add(4)
	b.Add(5)

	a.Add(b.ToSlice()...)

	if a.Len() != 5 {
		t.Error("set b is added with from an empty set a and therefore should have 5 elements in it")
	}

	d := NewUnsafeSet()
	d.Add(10)
	d.Add(14)
	d.Add(0)

	a.Add(d.ToSlice()...)
	if a.Len() != 8 {
		t.Error("set a should should have 8 elements in it after being added with with items from set d")
	}

	e := NewUnsafeSet()
	e.Add(14)
	e.Add(3)

	e.Add(a.ToSlice()...)
	if e.Len() != 8 {
		t.Error("set e should still have 8 elements in it after being added with items from set a that has duplicates")
	}
}

func TestSetIntersect(t *testing.T) {
	a := NewSet()
	a.Add(1)
	a.Add(3)
	a.Add(5)

	b := NewSet()
	a.Add(2)
	a.Add(4)
	a.Add(6)

	c := a.Intersect(b)

	if c.Len() != 0 {
		t.Error("set c should be the empty set because there is no common items to intersect")
	}

	a.Add(10)
	b.Add(10)

	d := a.Intersect(b)

	if !(d.Len() == 1 && d.Contains(10)) {
		t.Error("set d should have a size of 1 and contain the item 10")
	}
}

func TestUnsafeSetIntersect(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(1)
	a.Add(3)
	a.Add(5)

	b := NewUnsafeSet()
	a.Add(2)
	a.Add(4)
	a.Add(6)

	c := a.Intersect(b)

	if c.Len() != 0 {
		t.Error("set c should be the empty set because there is no common items to intersect")
	}

	a.Add(10)
	b.Add(10)

	d := a.Intersect(b)

	if !(d.Len() == 1 && d.Contains(10)) {
		t.Error("set d should have a size of 1 and contain the item 10")
	}
}

func TestSetDisjointFrom(t *testing.T) {
	a := NewSet()
	a.Add(1)
	a.Add(3)
	a.Add(5)

	b := NewSet()
	a.Add(2)
	a.Add(4)
	a.Add(6)

	if !a.IsDisjointFrom(b) {
		t.Error("both sets should be the disjoint from one another as there are no common items to intersect.")
	}

	a.Add(10)
	b.Add(10)

	if a.IsDisjointFrom(b) {
		t.Error("both sets a 1 item in common and should not be disjoint")
	}
}

func TestUnsafeSetDisjointFrom(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(1)
	a.Add(3)
	a.Add(5)

	b := NewUnsafeSet()
	a.Add(2)
	a.Add(4)
	a.Add(6)

	if !a.IsDisjointFrom(b) {
		t.Error("both sets should be the disjoint from one another as there are no common items to intersect.")
	}

	a.Add(10)
	b.Add(10)

	if a.IsDisjointFrom(b) {
		t.Error("both sets a 1 item in common and should not be disjoint")
	}
}

func TestSetDifference(t *testing.T) {
	a := NewSet()
	a.Add(1)
	a.Add(2)
	a.Add(3)

	b := NewSet()
	b.Add(1)
	b.Add(3)
	b.Add(4)
	b.Add(5)
	b.Add(6)
	b.Add(99)

	c := a.Difference(b)

	if !(c.Len() == 1 && c.Contains(2)) {
		t.Error("the difference of set a to b is the set of 1 item: 2")
	}
}

func TestUnsafeSetDifference(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(1)
	a.Add(2)
	a.Add(3)

	b := NewUnsafeSet()
	b.Add(1)
	b.Add(3)
	b.Add(4)
	b.Add(5)
	b.Add(6)
	b.Add(99)

	c := a.Difference(b)

	if !(c.Len() == 1 && c.Contains(2)) {
		t.Error("the difference of set a to b is the set of 1 item: 2")
	}
}

func TestSetEqual(t *testing.T) {
	a := NewSet()
	b := NewSet()

	if !a.Equal(b) {
		t.Error("Both a and b are empty sets, and should be equal")
	}

	a.Add(10)

	if a.Equal(b) {
		t.Error("a should not be equal to b because b is empty and a has item 1 in it")
	}

	b.Add(10)

	if !a.Equal(b) {
		t.Error("a is now equal again to b because both have the item 10 in them")
	}

	b.Add(8)
	b.Add(3)
	b.Add(47)

	if a.Equal(b) {
		t.Error("b has 3 more elements in it so therefore should not be equal to a")
	}

	a.Add(8)
	a.Add(3)
	a.Add(47)

	if !a.Equal(b) {
		t.Error("a and b should be equal with the same number of elements")
	}
}

func TestUnsafeSetEqual(t *testing.T) {
	a := NewUnsafeSet()
	b := NewUnsafeSet()

	if !a.Equal(b) {
		t.Error("Both a and b are empty sets, and should be equal")
	}

	a.Add(10)

	if a.Equal(b) {
		t.Error("a should not be equal to b because b is empty and a has item 1 in it")
	}

	b.Add(10)

	if !a.Equal(b) {
		t.Error("a is now equal again to b because both have the item 10 in them")
	}

	b.Add(8)
	b.Add(3)
	b.Add(47)

	if a.Equal(b) {
		t.Error("b has 3 more elements in it so therefore should not be equal to a")
	}

	a.Add(8)
	a.Add(3)
	a.Add(47)

	if !a.Equal(b) {
		t.Error("a and b should be equal with the same number of elements")
	}
}

func TestSetClone(t *testing.T) {
	a := NewSet()
	a.Add(1)
	a.Add(2)

	b := a.Clone()

	if !a.Equal(b) {
		t.Error("Clones should be equal")
	}

	a.Add(3)
	if a.Equal(b) {
		t.Error("a contains one more element, they should not be equal")
	}

	c := a.Clone()
	c.Remove(1)

	if a.Equal(c) {
		t.Error("C contains one element less, they should not be equal")
	}
}

func TestUnsafeSetClone(t *testing.T) {
	a := NewUnsafeSet()
	a.Add(1)
	a.Add(2)

	b := a.Clone()

	if !a.Equal(b) {
		t.Error("Clones should be equal")
	}

	a.Add(3)
	if a.Equal(b) {
		t.Error("a contains one more element, they should not be equal")
	}

	c := a.Clone()
	c.Remove(1)

	if a.Equal(c) {
		t.Error("C contains one element less, they should not be equal")
	}
}

func TestIter(t *testing.T) {
	a := NewSet()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := NewSet()
	for val := range a.Iter() {
		b.Add(val)
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iter) through the first set")
	}
}

func TestUnsafeIter(t *testing.T) {
	a := NewUnsafeSet()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := NewUnsafeSet()
	for val := range a.Iter() {
		b.Add(val)
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iter) through the first set")
	}
}

func TestIterator(t *testing.T) {
	a := NewSet()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := NewSet()
	it := a.Iterator()
	for it.HasNext() {
		b.Add(it.Next())
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iterator) through the first set")
	}
}

func TestUnsafeIterator(t *testing.T) {
	a := NewUnsafeSet()

	a.Add("Z")
	a.Add("Y")
	a.Add("X")
	a.Add("W")

	b := NewUnsafeSet()
	it := a.Iterator()
	for it.HasNext() {
		b.Add(it.Next())
	}

	if !a.Equal(b) {
		t.Error("The sets are not equal after iterating (Iterator) through the first set")
	}
}

func TestEmptySetProperties(t *testing.T) {
	empty := NewSet()

	a := NewSet()
	a.Add(1)
	a.Add("foo")
	a.Add("bar")

	b := NewSet()
	b.Add("one")
	b.Add("two")
	b.Add(3)
	b.Add(4)

	if !empty.IsSubset(a) || !empty.IsSubset(b) {
		t.Error("The empty set is supposed to be a subset of all sets")
	}

	if !a.IsSuperset(empty) || !b.IsSuperset(empty) {
		t.Error("All sets are supposed to be a superset of the empty set")
	}

	if !empty.IsSubset(empty) || !empty.IsSuperset(empty) {
		t.Error("The empty set is supposed to be a subset and a superset of itself")
	}

	c := a.Union(empty)
	if !c.Equal(a) {
		t.Error("The union of any set with the empty set is supposed to be equal to itself")
	}

	c = a.Intersect(empty)
	if !c.Equal(empty) {
		t.Error("The intesection of any set with the empty set is supposed to be the empty set")
	}

	if empty.Len() != 0 {
		t.Error("Len of the empty set is supposed to be zero")
	}

}

func TestToSliceUnthreadsafe(t *testing.T) {
	s := makeUnsafeSet([]int{1, 2, 3})
	setAsSlice := s.ToSlice()
	if len(setAsSlice) != s.Len() {
		t.Errorf("Set length is incorrect: %v", len(setAsSlice))
	}

	for _, i := range setAsSlice {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}
