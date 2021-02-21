package tx

import (
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
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
