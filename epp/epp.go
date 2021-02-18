package epp

import (
	"github.com/jtejido/ngac/internal/set"
	"github.com/jtejido/ngac/pap"
	"github.com/jtejido/ngac/pdp"
	"github.com/jtejido/ngac/pip/graph"
	"github.com/jtejido/ngac/pip/obligations"
)

type EPP struct {
	pap               *pap.PAP
	pdp               *pdp.PDP
	functionEvaluator *FunctionEvaluator
}

func NewEPP(p *pap.PDP, eppOptions *EPPOptions) *EPP {
	e := new(EPP)
	e.pap = p.PAP()
	e.pdp = p
	e.functionEvaluator = NewFunctionEvaluator()
	if eppOptions != nil {
		for _, executor := range eppOptions.Executors() {
			e.functionEvaluator.Add(executor)
		}
	}

	return e
}

func (epp *EPP) ProcessEvent(eventCtx EventContext, user, process string) {
	obligations := epp.pap.ObligationsPAP().All()
	for obligation := range obligations.Iter() {
		if !obligation.(*obligations.Obligation).Enabled {
			continue
		}

		for _, rule := range obligation.Rules {
			if !epp.eventMatches(user, process, eventCtx.Event(), eventCtx.Target(), rule.EventPattern) {
				continue
			}

			// check the response condition
			responsePattern := rule.ResponsePattern
			if !pr.checkCondition(responsePattern.Condition, eventCtx, user, process, epp.pdp) {
				continue
			} else if !pr.checkNegatedCondition(responsePattern.NegatedCondition, eventCtx, user, process, epp.pdp) {
				continue
			}

			for _, action := range rule.ResponsePattern.Actions {
				if !pr.checkCondition(action.Condition, eventCtx, user, process, epp.pdp) {
					continue
				} else if !pr.checkNegatedCondition(action.NegatedCondition, eventCtx, user, process, epp.pdp) {
					continue
				}

				pr.applyAction(obligation.getLabel(), eventCtx, user, process, action)
			}
		}
	}
}

func (epp *EPP) checkCondition(condition *obligations.Condition, eventCtx EventContext, user, process string) bool {
	if condition == nil {
		return true
	}

	functions := condition.Condition
	for _, f := range condition.Condition {
		result := epp.functionEvaluator.evalBool(eventCtx, user, process, epp.pdp, f)
		if !result {
			return false
		}
	}

	return true
}

func (epp *EPP) checkNegatedCondition(condition *obligations.NegatedCondition, eventCtx EventContext, user, process string) bool {
	if condition == nil {
		return true
	}

	functions := condition.Condition
	for _, f := range condition.Condition {
		result := epp.functionEvaluator.evalBool(eventCtx, user, process, epp.pdp, f)
		if result {
			return false
		}
	}

	return true
}

func (epp *EPP) eventMatches(user, process, event string, target *graph.Node, match *obligations.EventPattern) bool {
	if match.Operations != nil {
		var found bool
		for _, v := range match.Operations {
			if v == event {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	matchSubject := match.Subject
	matchPolicyClass := match.PolicyClass
	matchTarget := match.Target

	return epp.subjectMatches(user, process, matchSubject) && epp.pcMatches(user, matchPolicyClass) && epp.targetMatches(target, matchTarget)
}

func (epp *EPP) subjectMatches(user, process string, matchSubject *obligations.Subject) bool {
	if matchSubject == nil {
		return true
	}

	// any user
	if (matchSubject.AnyUser == nil && matchSubject.Process == nil) || (matchSubject.AnyUser != nil && len(matchSubject.AnyUser) == 0) {
		return true
	}

	// get the current user node
	userNode, err := epp.pap.Graph.Node(user)
	if err != nil {
		panic(err.Error())
	}

	if epp.checkAnyUser(userNode, matchSubject.AnyUser) {
		return true
	}

	if matchSubject.User == userNode.Name() {
		return true
	}

	return matchSubject.Process != nil && matchSubject.Process.Value == process
}

func (epp *EPP) checkAnyUser(userNode *graph.Node, anyUser []string) bool {
	if anyUser == nil || len(anyUser) == 0 {
		return true
	}

	dfs := graph.NewDFS(epp.pdp.PAP().Graph)

	// check each user in the anyUser list
	// there can be users and user attributes
	for _, u := range anyUser {
		anyUserNode, err := epp.pdp.PAP().Graph.Node(u)
		if err != nil {
			panic(err.Error())
		}

		// if the node in anyUser == the user than return true
		if anyUserNode.Name == userNode.Name {
			return true
		}

		// if the anyUser is not an UA, move to the next one
		if anyUserNode.Type != graph.UA {
			continue
		}

		nodes := set.NewSet()
		visitor := func(node *graph.Node) error {
			if node.Name == userNode.Name {
				nodes.Add(node.Name)
			}
		}
		propagator := func(parent, child *graph.Node) error { return nil }

		err = dfs.Traverse(userNode, propagator, visitor, graph.PARENTS)
		if err != nil {
			return panic(err.Error())
		}

		if nodes.Contains(anyUserNode.Name) {
			return true
		}
	}

	return false
}

func (epp *EPP) pcMatches(user string, matchPolicyClass *obligations.PolicyClass) bool {
	// not yet implemented
	return true
}

func (epp *EPP) targetMatches(target *graph.Node, matchTarget *obligations.Target) bool {
	if matchTarget == nil {
		return true
	}

	if matchTarget.PolicyElements == nil && matchTarget.Containers == nil {
		return true
	}

	if matchTarget.Containers != nil {
		if len(matchTarget.Containers) == 0 {
			return true
		}

		// check that target is contained in any container
		containers := epp.containersOf(target.Name)
		for _, evrContainer := range matchTarget.Containers {
			it := containers.Iterator()
			for it.HasNext() {
				contNode := it.Next().(*graph.Node)
				if epp.nodesMatch(evrContainer, contNode) {
					return true
				}
			}
		}

		return false
	} else if matchTarget.PolicyElements != nil {
		if len(matchTarget.PolicyElements) == 0 {
			return true
		}

		// check that target is in the list of policy elements
		for _, evrNode := range matchTarget.PolicyElements {
			if epp.nodesMatch(evrNode, target) {
				return true
			}
		}

		return false
	}

	return false
}
