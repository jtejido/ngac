package pap

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
	"github.com/jtejido/ngac/pip/tx"
)

type PAP struct {
	graphAdmin        *GraphAdmin
	prohibitionsAdmin *ProhibitionsAdmin
	obligationsAdmin  *ObligationsAdmin
}

func NewPAP(p common.FunctionalEntity) *PAP {
	return &PAP{NewGraphAdmin(p), NewProhibitionsAdmin(p), NewObligationsAdmin(p)}
}

func (pap *PAP) Graph() graph.Graph {
	return pap.graphAdmin
}

func (pap *PAP) Prohibitions() prohibitions.Prohibitions {
	return pap.prohibitionsAdmin
}

func (pap *PAP) Obligations() obligations.Obligations {
	return pap.obligationsAdmin
}

func (pap *PAP) RunTx(txRunner common.TxRunner) error {
	tx := tx.NewMemTx(pap.graphAdmin, pap.prohibitionsAdmin, pap.obligationsAdmin)
	return tx.RunTx(txRunner)
}
