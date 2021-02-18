package common

import (
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
)

type TxRunner func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error

type FunctionalEntity interface {
	Graph() graph.Graph
	Prohibitions() prohibitions.Prohibitions
	Obligations() obligations.Obligations
	RunTx(txRunner TxRunner) error
}
