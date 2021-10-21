package pap

import (
	"ngac/pkg/common"
	"ngac/pkg/pip/obligations"
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

func (oa *ObligationsAdmin) Add(obligation *obligations.Obligation, enable bool) {
	oa.obligations.Add(obligation, enable)
}

func (oa *ObligationsAdmin) Get(label string) *obligations.Obligation {
	return oa.obligations.Get(label)
}

func (oa *ObligationsAdmin) All() []*obligations.Obligation {
	return oa.obligations.All()
}

func (oa *ObligationsAdmin) Update(label string, obligation *obligations.Obligation) {
	oa.obligations.Update(label, obligation)
}

func (oa *ObligationsAdmin) Remove(label string) {
	oa.obligations.Remove(label)
}

func (oa *ObligationsAdmin) SetEnable(label string, enabled bool) {
	oa.obligations.SetEnable(label, enabled)
}

func (oa *ObligationsAdmin) GetEnabled() []*obligations.Obligation {
	return oa.obligations.GetEnabled()
}
