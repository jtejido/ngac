package common

import (
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

type TxRunner func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error

type PolicyStore interface {
	Graph() graph.Graph
	Prohibitions() prohibitions.Prohibitions
	Obligations() obligations.Obligations
	RunTx(txRunner TxRunner) error
}
