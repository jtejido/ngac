package common

import (
	"github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/jtejido/ngac/pkg/pip/obligations"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type TxRunner func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error

type PolicyStore interface {
	Graph() graph.Graph
	Prohibitions() prohibitions.Prohibitions
	Obligations() obligations.Obligations
	RunTx(txRunner TxRunner) error
}
