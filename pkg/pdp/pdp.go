package pdp

import (
	"ngac/pkg/common"
	"ngac/pkg/context"
	"ngac/pkg/epp"
	"ngac/pkg/pdp/audit"
	"ngac/pkg/pdp/decider"
	"ngac/pkg/pdp/service"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
	"ngac/pkg/pip/tx"
)

type PDP struct {
	epp     *EPP
	pap     common.FunctionalEntity
	decider decider.Decider
	auditor audit.Auditor
}

// func NewPDP(pap *pap.PAP, eppOptions *epp.EPPOptions) (pdp *PDP, err error) {
func NewPDP(pap common.FunctionalEntity, eppOptions *epp.EPPOptions, decider decider.Decider, auditor audit.Auditor) *PDP {
	// create PDP
	pdp := newPDP(pap, decider, auditor)
	// create the EPP
	pdp.epp = NewEPP(pap, pdp, eppOptions)
	pdp.decider = decider
	pdp.auditor = auditor

	return pdp
}

/**
 * Create a new PDP instance given a Policy Administration Point and an optional set of FunctionExecutors to be
 * used by the EPP.
 * @param pap the Policy Administration Point that the PDP will use to change the graph.
 * @throws PMException if there is an error initializing the EPP.
 */
func newPDP(pap common.FunctionalEntity, decider decider.Decider, auditor audit.Auditor) *PDP {
	return &PDP{pap: pap, decider: decider, auditor: auditor}
}

type WithUser struct {
	userCtx context.Context
	pap     common.FunctionalEntity
	epp     *EPP
	decider decider.Decider
	auditor audit.Auditor
	gs      graph.Graph
	ps      prohibitions.Prohibitions
	os      obligations.Obligations
	// analyticsService
}

var _ common.FunctionalEntity = &WithUser{}

func (p *PDP) WithUser(userCtx context.Context) *WithUser {
	return newWithUser(userCtx, p.pap, p.epp, p.decider, p.auditor)
}

func newWithUser(u context.Context, p common.FunctionalEntity, e *EPP, d decider.Decider, a audit.Auditor) *WithUser {
	return &WithUser{u, p, e, d, a, service.NewGraphService(u, p, e, d, a), service.NewProhibitionsService(u, p, e, d, a), service.NewObligationsService(u, p, e, d, a)}
}

func (wu *WithUser) Graph() graph.Graph {
	return wu.gs
}

func (wu *WithUser) Prohibitions() prohibitions.Prohibitions {
	return wu.ps
}

func (wu *WithUser) Obligations() obligations.Obligations {
	return wu.os
}

func (wu *WithUser) RunTx(txRunner common.TxRunner) error {
	tx := tx.NewMemTx(wu.gs, wu.ps, wu.os)
	return tx.RunTx(txRunner)
}
