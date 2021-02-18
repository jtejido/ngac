package set

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const N = 1000

func TestAddConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
	for _, i := range ints {
		if !s.Contains(i) {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}

func TestLenConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		elems := s.Len()
		for i := 0; i < N; i++ {
			newElems := s.Len()
			if newElems < elems {
				t.Errorf("Len shrunk from %v to %v", elems, newElems)
			}
		}
		wg.Done()
	}()

	for i := 0; i < N; i++ {
		s.Add(rand.Int())
	}
	wg.Wait()
}

func TestClearConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func() {
			s.Clear()
			wg.Done()
		}()
		go func(i int) {
			s.Add(i)
		}(i)
	}

	wg.Wait()
}

func TestCloneConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)

	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := range ints {
		go func(i int) {
			s.Remove(i)
			wg.Done()
		}(i)
	}

	s.Clone()
}

func TestContainsConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)
	interfaces := make([]interface{}, 0)
	for _, v := range ints {
		s.Add(v)
		interfaces = append(interfaces, v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.Contains(interfaces...)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestDifferenceConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.Difference(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestAddFromConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.AddFrom(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestEqualConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.Equal(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIntersectConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.Intersect(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIsSubsetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.IsSubset(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIsProperSubsetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.IsProperSubset(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIsSupersetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.IsSuperset(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIsProperSupersetConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s, ss := NewSet(), NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
		ss.Add(v)
	}

	var wg sync.WaitGroup
	for range ints {
		wg.Add(1)
		go func() {
			s.IsProperSuperset(ss)
			wg.Done()
		}()
	}
	wg.Wait()
}

func TestIterConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
	}

	cs := make([]<-chan interface{}, 0)
	for range ints {
		cs = append(cs, s.Iter())
	}

	c := make(chan interface{})
	go func() {
		for n := 0; n < len(ints)*N; {
			for _, d := range cs {
				select {
				case <-d:
					n++
					c <- nil
				default:
				}
			}
		}
		close(c)
	}()

	for range c {
	}
}

func TestRemoveConcurrent(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)
	for _, v := range ints {
		s.Add(v)
	}

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for _, v := range ints {
		go func(i int) {
			s.Remove(i)
			wg.Done()
		}(v)
	}
	wg.Wait()

	if s.Len() != 0 {
		t.Errorf("Expected Len 0; got %v", s.Len())
	}
}

func TestToSlice(t *testing.T) {
	runtime.GOMAXPROCS(2)

	s := NewSet()
	ints := rand.Perm(N)

	var wg sync.WaitGroup
	wg.Add(len(ints))
	for i := 0; i < len(ints); i++ {
		go func(i int) {
			s.Add(i)
			wg.Done()
		}(i)
	}

	wg.Wait()
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

func TestToSliceDeadlock(t *testing.T) {
	runtime.GOMAXPROCS(2)

	var wg sync.WaitGroup
	set := NewSet()
	workers := 10
	wg.Add(workers)
	for i := 1; i <= workers; i++ {
		go func() {
			for j := 0; j < 1000; j++ {
				set.Add(1)
				set.ToSlice()
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
