package obligations

import (
	"math/rand"
	"runtime"
	"sync"
	"testing"
)

const (
	N      = 1000
	strLen = 10
)

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func TestAddGetConcurrent(t *testing.T) {
	pros := make([]string, N)
	for i := 0; i < N; i++ {
		pros[i] = randomString(strLen)
	}
	runtime.GOMAXPROCS(2)

	s := NewMemObligations()
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			ob := NewObligation(pros[i])
			ob.Label = pros[i]
			s.Add(ob, false)
			wg.Done()
		}(i)
	}

	wg.Wait()
	for i := 0; i < N; i++ {
		if s.Get(pros[i]) == nil {
			t.Errorf("Set is missing element: %v", i)
		}
	}
}
