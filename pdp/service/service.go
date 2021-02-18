package service

import (
	"github.com/jtejido/ngac/pdp/decider"
	"github.com/jtejido/ngac/pdp/policy"
	"github.com/jtejido/ngac/pip/graph"
)

type Service struct {
	decider     decider.Decider
	superPolicy *policy.SuperPolicy
}

func (s *Service) Decider() decider.Decider {
	return s.decider
}

func (s *Service) hasPermissions(ctx Context, targetNode *graph.Node, permissions ...interface{}) bool {
	name := targetNode.Name
	if targetNode.Type == graph.PC {
		temp, found := targetNode.Properties.Get(graph.REP_PROPERTY)
		if !found {
			return false
		}
		name = temp.(string)
	}

	perms := s.decider.List(ctx.User(), ctx.Process(), name)

	for _, p := range permissions {
		if p == decider.ANY_OPERATIONS {
			return perms.Len() > 0
		}
	}

	if perms.Contains(graph.ALL_OPS) {
		return true
	}

	return perms.Len() > 0 || perms.Contains(permissions...)

}
