package pip

import (
    "github.com/jtejido/ngac/pkg/common"
    "github.com/jtejido/ngac/pkg/pip/graph"
    "github.com/jtejido/ngac/pkg/pip/obligations"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
    "github.com/jtejido/ngac/pkg/pip/tx"
)

var _ common.PolicyStore = &PIP{}

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
