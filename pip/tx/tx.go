package tx

import (
	"github.com/jtejido/ngac/pip/Prohibitions"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"sync"
)

type Tx struct {
	sync.RWMutex
	graph        graph.Graph
	prohibitions prohibitions.Prohibitions
	obligations  obligations.Obligations
}
