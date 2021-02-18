package epp

import (
	"github.com/jtejido/ngac/common"
	"github.com/jtejido/ngac/pdp"
	"github.com/jtejido/ngac/pdp/service"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
	"github.com/jtejido/ngac/pip/prohibitions"
)

type EPP struct {
	pap               common.FunctionalEntity
	pdp               *pdp.PDP
	functionEvaluator *FunctionEvaluator
}

func NewEPP(pap common.FunctionalEntity, p *pdp.PDP, eppOptions *EPPOptions) *EPP {
	e := new(EPP)
	e.pap = pap
	e.pdp = p
	e.functionEvaluator = NewFunctionEvaluator()
	if eppOptions != nil {
		for _, executor := range eppOptions.Executors() {
			e.functionEvaluator.Add(executor)
		}
	}

	return e
}

func (epp *EPP) AddFunctionExecutor(executor FunctionExecutor) {
	epp.functionEvaluator.Add(executor)
}

func (epp *EPP) RemoveFunctionExecutor(executor FunctionExecutor) {
	epp.functionEvaluator.Remove(executor)
}

func (epp *EPP) ProcessEvent(eventCtx EventContext, user, process string) error {
	obligs := epp.pap.Obligations().All()
	for _, obligation := range obligs {
		if !obligation.Enabled {
			continue
		}

		definingUser, _ := service.NewUserContext(obligation.User)

		if err := epp.pdp.WithUser(definingUser).RunTx(func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error {
			rules := obligation.Rules
			for _, rule := range rules {
				if !eventCtx.MatchesPattern(rule.EventPattern, g) {
					continue
				}

				responsePattern := rule.ResponsePattern
				responsePattern.Apply(g, p, o, epp.functionEvaluator, eventCtx, rule, obligation.Label)
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
