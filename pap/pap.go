package pap

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
	"github.com/jtejido/ngac/pip/tx"
)

var _ common.FunctionalEntity = &PAP{}

type PAP struct {
	graphAdmin        *GraphAdmin
	prohibitionsAdmin *ProhibitionsAdmin
	obligationsAdmin  *ObligationsAdmin
}

func NewPAP(p common.FunctionalEntity) (*PAP, error) {
	ga, err := NewGraphAdmin(p)
	if err != nil {
		return nil, err
	}
	return &PAP{ga, NewProhibitionsAdmin(p), NewObligationsAdmin(p)}, nil
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
