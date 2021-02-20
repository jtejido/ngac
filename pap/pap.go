package pap

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pip/tx"
)

type PAP struct {
	graphAdmin        GraphAdmin
	prohibitionsAdmin ProhibitionsAdmin
	obligationsAdmin  ObligationsAdmin
}

func NewPAP(graphAdmin GraphAdmin, prohibitionsAdmin ProhibitionsAdmin, obligationsAdmin ObligationsAdmin) *PAP {
	return &PAP{graphAdmin, prohibitionsAdmin, obligationsAdmin}
}

func (pap *PAP) Graph() Graph {
	return pap.graphAdmin
}

func (pap *PAP) Prohibitions() Prohibitions {
	return pap.prohibitionsAdmin
}

func (pap *PAP) Obligations() Obligations {
	return pap.obligationsAdmin
}

func (pap *PAP) RunTx(txRunner common.TxRunner) error {
	tx := tx.NewMemTx(pap.graphAdmin, pap.prohibitionsAdmin, pap.obligationsAdmin)
	return tx.RunTx(txRunner)
}
