package service

import (
    "github.com/jtejido/ngac/audit"
    "github.com/jtejido/ngac/common"
    "github.com/jtejido/ngac/context"
    "github.com/jtejido/ngac/decider"
    "github.com/jtejido/ngac/epp"
    "github.com/jtejido/ngac/guard"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/pip/obligations"
)

var _ obligations.Obligations = &Obligations{}

type Obligations struct {
    Service
    guard *guard.Obligations
}

func NewObligationsService(userCtx context.Context, p common.FunctionalEntity, e epp.EPP, d decider.Decider, a audit.Auditor) *Obligations {
    ans := new(Obligations)
    ans.userCtx = userCtx
    ans.pap = p
    ans.epp = e
    ans.decider = d
    ans.auditor = a
    ans.guard = guard.NewObligationsGuard(p, d)
    return ans
}

func (o *Obligations) Add(obligation *obligations.Obligation, enable bool) {
    o.guard.CheckAdd(o.userCtx)
    o.ObligationsAdmin().Add(obligation, enable)
}

func (o *Obligations) Get(label string) *obligations.Obligation {
    o.guard.CheckGet(o.userCtx)
    return o.ObligationsAdmin().Get(label)
}

func (o *Obligations) All() []*obligations.Obligation {
    o.guard.CheckGet(o.userCtx)
    return o.ObligationsAdmin().All()
}

func (o *Obligations) Update(label string, obligation *obligations.Obligation) {
    o.guard.CheckUpdate(o.userCtx)
    o.ObligationsAdmin().Update(label, obligation)
}

func (o *Obligations) Remove(label string) {
    o.guard.CheckDelete(o.userCtx)
    o.ObligationsAdmin().Remove(label)
}

func (o *Obligations) SetEnable(label string, enabled bool) {
    o.guard.CheckEnable(o.userCtx)
    o.ObligationsAdmin().SetEnable(label, enabled)
}

func (o *Obligations) GetEnabled() []*obligations.Obligation {
    o.guard.CheckGet(o.userCtx)
    return o.ObligationsAdmin().GetEnabled()
}

func (o *Obligations) Reset() (err error) {
    o.guard.CheckReset(o.userCtx)

    obs := o.All()
    labels := set.NewSet()
    for _, obli := range obs {
        labels.Add(obli.Label)
    }
    for label := range labels.Iter() {
        o.Remove(label.(string))
    }

    return nil
}
