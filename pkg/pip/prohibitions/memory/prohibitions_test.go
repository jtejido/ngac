package memory

import (
	"math/rand"
	"ngac/pkg/operations"
	p "ngac/pkg/pip/prohibitions"
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

	s := New()
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			s.Add(p.NewProhibition(pros[i], pros[i], nil, nil, false))
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

func TestCreateProhibition(t *testing.T) {
	prohibs := New()

	builder := p.NewBuilder("prohibition1", "123", operations.NewOperationSet("read"))
	builder.AddContainer("1234", true)
	prohibition := builder.Build()
	prohibs.Add(prohibition)

	builder = p.NewBuilder("p123", "sub", operations.NewOperationSet("read"))
	builder.AddContainer("1234", true)
	prohibition = builder.Build()

	prohibs.Add(prohibition)

	p123 := prohibs.Get("p123")
	if p123.Name != "p123" {
		t.Errorf("Name do not match")
	}
	if p123.Subject != "sub" {
		t.Errorf("Subject do not match")
	}
	if p123.Intersection {
		t.Errorf("Intersection should be false")
	}
	if p123.Intersection {
		t.Errorf("Intersection should be false")
	}
	if _, ok := prohibition.Containers()["1234"]; !ok {
		t.Errorf("\\'1234\\' should be in containers")
	}
	if v, ok := prohibition.Containers()["1234"]; !ok {
		t.Errorf("\\'1234\\' should be in containers")
	} else {
		if !v {
			t.Errorf("\\'1234\\' should be \\'true\\'")
		}
	}

}

func TestGetProhibitions(t *testing.T) {
	prohibs := New()

	builder := p.NewBuilder("prohibition1", "123", operations.NewOperationSet("read"))
	builder.AddContainer("1234", true)
	prohibition := builder.Build()
	prohibs.Add(prohibition)

	prohibitions := prohibs.All()
	if len(prohibitions) != 1 {
		t.Errorf("incorrect size")
	}
}

func TestGetProhibition(t *testing.T) {
	prohibs := New()

	builder := p.NewBuilder("prohibition1", "123", operations.NewOperationSet("read"))
	builder.AddContainer("1234", true)
	prohibition := builder.Build()
	prohibs.Add(prohibition)

	prohibition = prohibs.Get("prohibition1")
	if prohibition.Name != "prohibition1" {
		t.Errorf("incorrect name")
	}
	if prohibition.Subject != "123" {
		t.Errorf("incorrect subject")
	}
	if prohibition.Intersection {
		t.Errorf("Intersection should be false")
	}
	if _, ok := prohibition.Containers()["1234"]; !ok {
		t.Errorf("\\'1234\\' should be in containers")
	}
	if v, ok := prohibition.Containers()["1234"]; !ok {
		t.Errorf("\\'1234\\' should be in containers")
	} else {
		if !v {
			t.Errorf("\\'1234\\' should be \\'true\\'")
		}
	}

}
