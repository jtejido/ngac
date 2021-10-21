package tx

import (
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
	"sync"
)

type Committer func() error

type Tx struct {
	sync.RWMutex
	graph        graph.Graph
	prohibitions prohibitions.Prohibitions
	obligations  obligations.Obligations
}

func assert(t bool) {
	if !t {
		panic("assertion failed")
	}
}
