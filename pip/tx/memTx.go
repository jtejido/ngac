package tx

import (
    "github.com/jtejido/ngac/common"
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type MemTx struct {
    Tx
    txGraph        *TxGraph
    txProhibitions *TxProhibitions
    txObligations  *TxObligations
}

func NewMemTx(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) *MemTx {
    ans := new(MemTx)
    ans.graph = g
    ans.prohibitions = p
    ans.obligations = o
    ans.txGraph = NewTxGraph(g)
    ans.txProhibitions = NewTxProhibitions(p)
    ans.txObligations = NewTxObligations(o)
    return ans
}

func (mt *MemTx) RunTx(txRunner common.TxRunner) error {
    if err := txRunner.Run(txGraph, txProhibitions, txObligations); err != nil {
        mt.Rollback()
        return err
    }
    mt.Commit()
    return nil
}

func (mt *MemTx) Commit() error {
    ans.Lock()
    defer ans.Unlock()
    // commit the graph
    txGraph.Commit()

    // commit the prohibitions
    txProhibitions.Commit()

    // commit the obligations
    txObligations.Commit()
}

func (mt *MemTx) Rollback() {
    // rollback graph
    mt.txGraph = NewTxGraph(mt.graph)

    // rollback prohibitions
    mt.txProhibitions = NewTxProhibitions(mt.prohibitions)

    // rollback obligations
    mt.txObligations = NewTxObligations(mt.obligations)
}
