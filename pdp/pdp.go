package pdp

import (
	// "github.com/jtejido/ngac/epp"
	"github.com/jtejido/ngac/pap"
	"github.com/jtejido/ngac/pdp/audit"
	"github.com/jtejido/ngac/pdp/decider"
)

type TokenString string

type PDP struct {
	// epp                 *epp.EPP
	pap     *pap.PAP
	decider decider.Decider
	auditor audit.Auditor
}

// func NewPDP(pap *pap.PAP, eppOptions *epp.EPPOptions) (pdp *PDP, err error) {
func NewPDP(pap *pap.PAP, decider decider.Decider, auditor audit.Auditor) (pdp *PDP, err error) {
	// pdp.epp = epp.NewEPP(pdp, eppOptions)

	pdp.pap = pap
	pdp.decider = decider
	pdp.auditor = auditor

	return pdp, nil
}
