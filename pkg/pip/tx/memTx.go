package tx

import (
    "github.com/jtejido/ngac/pkg/common"
    "github.com/jtejido/ngac/pkg/pip/graph"
    "github.com/jtejido/ngac/pkg/pip/obligations"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
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
    // should be clone to revert back to its original state
    ans.txGraph = NewTxGraph(g)
    ans.txProhibitions = NewTxProhibitions(p)
    ans.txObligations = NewTxObligations(o)
    return ans
}

func (mt *MemTx) RunTx(txRunner common.TxRunner) error {
    if err := txRunner(mt.txGraph, mt.txProhibitions, mt.txObligations); err != nil {
        mt.Rollback()
        return err
    }
    return mt.Commit()
}

func (mt *MemTx) Commit() (err error) {
    mt.Lock()
    defer mt.Unlock()
    // commit the graph
    if err = mt.txGraph.Commit(); err == nil {
        // commit the prohibitions
        if err = mt.txProhibitions.Commit(); err == nil {
            // commit the obligations
            if err = mt.txObligations.Commit(); err == nil {
                return
            }
        }
    }

    return
}

func (mt *MemTx) Rollback() {
    // rollback graph
    mt.txGraph = NewTxGraph(mt.graph)

    // rollback prohibitions
    mt.txProhibitions = NewTxProhibitions(mt.prohibitions)

    // rollback obligations
    mt.txObligations = NewTxObligations(mt.obligations)
}
