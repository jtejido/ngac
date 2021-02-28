package service

import (
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/pkg/common"
    "github.com/jtejido/ngac/pkg/context"
    "github.com/jtejido/ngac/pkg/epp"
    "github.com/jtejido/ngac/pkg/pdp/audit"
    "github.com/jtejido/ngac/pkg/pdp/decider"
    "github.com/jtejido/ngac/pkg/pdp/service/guard"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
)

var _ prohibitions.Prohibitions = &Prohibitions{}

type Prohibitions struct {
    Service
    guard *guard.Prohibitions
}

func NewProhibitionsService(userCtx context.Context, p common.FunctionalEntity, e epp.EPP, d decider.Decider, a audit.Auditor) *Prohibitions {
    ans := new(Prohibitions)
    ans.userCtx = userCtx
    ans.pap = p
    ans.epp = e
    ans.decider = d
    ans.auditor = a
    ans.guard = guard.NewProhibitionsGuard(p, d)
    return ans
}

func (p *Prohibitions) Add(prohibition *prohibitions.Prohibition) {
    p.guard.CheckAdd(p.userCtx, prohibition)

    //create prohibition in PAP
    p.ProhibitionsAdmin().Add(prohibition)
}

func (p *Prohibitions) All() []*prohibitions.Prohibition {
    all := p.ProhibitionsAdmin().All()
    p.guard.Filter(p.userCtx, all)
    return all
}

func (p *Prohibitions) Get(prohibitionName string) *prohibitions.Prohibition {
    prohibition := p.ProhibitionsAdmin().Get(prohibitionName)
    p.guard.CheckGet(p.userCtx, prohibition)

    return prohibition
}

func (p *Prohibitions) ProhibitionsFor(subject string) []*prohibitions.Prohibition {
    prohibitionsFor := p.ProhibitionsAdmin().ProhibitionsFor(subject)
    p.guard.Filter(p.userCtx, prohibitionsFor)
    return prohibitionsFor
}

func (p *Prohibitions) Update(prohibitionName string, prohibition *prohibitions.Prohibition) {
    p.guard.CheckUpdate(p.userCtx, prohibition)
    p.ProhibitionsAdmin().Update(prohibitionName, prohibition)
}

func (p *Prohibitions) Remove(prohibitionName string) {
    p.guard.CheckDelete(p.userCtx, p.ProhibitionsAdmin().Get(prohibitionName))
    p.ProhibitionsAdmin().Remove(prohibitionName)
}

func (p *Prohibitions) Reset(userCtx context.Context) (err error) {
    p.guard.CheckReset(p.userCtx)

    prohs := p.ProhibitionsAdmin().All()
    names := set.NewSet()
    for _, pro := range prohs {
        names.Add(pro.Name)
    }
    for name := range names.Iter() {
        p.Remove(name.(string))
    }

    return nil
}
