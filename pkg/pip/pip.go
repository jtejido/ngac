package pip

import (
    "ngac/pkg/common"
    "ngac/pkg/pip/graph"
    "ngac/pkg/pip/obligations"
    "ngac/pkg/pip/prohibitions"
    "ngac/pkg/pip/tx"
)

var _ common.FunctionalEntity = &PIP{}

type PIP struct {
    graph        graph.Graph
    prohibitions prohibitions.Prohibitions
    obligations  obligations.Obligations
}

func NewPIP(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) *PIP {
    return &PIP{g, p, o}
}

func (p *PIP) Graph() graph.Graph {
    return p.graph
}

func (p *PIP) Prohibitions() prohibitions.Prohibitions {
    return p.prohibitions
}

func (p *PIP) Obligations() obligations.Obligations {
    return p.obligations
}

func (p *PIP) RunTx(txRunner common.TxRunner) error {
    tx := tx.NewMemTx(p.graph, p.prohibitions, p.obligations)
    return tx.RunTx(txRunner)
}
