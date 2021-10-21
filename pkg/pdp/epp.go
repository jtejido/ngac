package pdp

import (
	"ngac/pkg/common"
	"ngac/pkg/context"
	"ngac/pkg/epp"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

type EPP struct {
	pap               common.FunctionalEntity
	pdp               *PDP
	functionEvaluator *epp.FunctionEvaluator
}

func NewEPP(pap common.FunctionalEntity, p *PDP, eppOptions *epp.EPPOptions) *EPP {
	e := new(EPP)
	e.pap = pap
	e.pdp = p
	e.functionEvaluator = epp.NewFunctionEvaluator()
	if eppOptions != nil {
		for _, executor := range eppOptions.Executors() {
			e.functionEvaluator.Add(executor)
		}
	}

	return e
}

func (e *EPP) AddFunctionExecutor(executor epp.FunctionExecutor) {
	e.functionEvaluator.Add(executor)
}

func (e *EPP) RemoveFunctionExecutor(executor epp.FunctionExecutor) {
	e.functionEvaluator.Remove(executor)
}

func (e *EPP) ProcessEvent(eventCtx epp.EventContext) error {
	obligs := e.pap.Obligations().All()
	for _, obligation := range obligs {
		if !obligation.Enabled {
			continue
		}

		definingUser, _ := context.NewUserContext(obligation.User)

		if err := e.pdp.WithUser(definingUser).RunTx(func(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations) error {
			rules := obligation.Rules
			for _, rule := range rules {
				if !eventCtx.MatchesPattern(rule.EventPattern, g) {
					continue
				}

				err := epp.Apply(g, p, o, e.functionEvaluator, eventCtx, rule, obligation.Label)
				if err != nil {
					return err
				}
			}

			return nil
		}); err != nil {
			return err
		}
	}

	return nil
}
