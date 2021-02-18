package pap

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pip/obligations"
)

var (
	_ obligations.Obligations = &ObligationsAdmin{}
)

type ObligationsAdmin struct {
	obligations obligations.Obligations
}

func NewObligationsAdmin(pip common.FunctionalEntity) *ObligationsAdmin {
	return &ObligationsAdmin{pip.Obligations()}
}

func (oa *ObligationsAdmin) AddObligation(obligation *obligations.Obligation, enable bool) error {
	return oa.obligations.AddObligation(obligation, enable)
}

func (oa *ObligationsAdmin) GetObligation(label string) *obligations.Obligation {
	return oa.obligations.GetObligation(label)
}

func (oa *ObligationsAdmin) All() []*obligations.Obligation {
	return oa.obligations.All()
}

func (oa *ObligationsAdmin) UpdateObligation(label string, obligation *obligations.Obligation) {
	oa.obligations.UpdateObligation(label, obligation)
}

func (oa *ObligationsAdmin) RemoveObligation(label string) {
	oa.obligations.RemoveObligation(label)
}

func (oa *ObligationsAdmin) SetEnable(label string, enabled bool) {
	oa.obligations.SetEnable(label, enabled)
}

func (oa *ObligationsAdmin) GetEnabled() []*obligations.Obligation {
	return oa.obligations.GetEnabled()
}
