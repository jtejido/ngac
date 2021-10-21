package epp

import (
	"fmt"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

var _ FunctionExecutor = &ChildOfAssignExecutor{}

type ChildOfAssignExecutor struct{}

func (f *ChildOfAssignExecutor) Name() string {
	return "child_of_assign"
}
func (f *ChildOfAssignExecutor) NumParams() int {
	return 0
}
func (f *ChildOfAssignExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
	eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
	var child *graph.Node
	if ev, ok := eventCtx.(*AssignToEvent); ok {
		child = ev.ChildNode
	} else if _, ok := eventCtx.(*AssignEvent); ok {
		child = eventCtx.Target()
	} else if ev, ok := eventCtx.(*DeassignFromEvent); ok {
		child = ev.ChildNode
	} else if _, ok := eventCtx.(*DeassignEvent); ok {
		child = eventCtx.Target()
	} else {
		return nil, fmt.Errorf("invalid event context for function child_of_assign. Valid event contexts are AssignTo, Assign, DeassignFrom, and Deassign")
	}

	return child, nil
}
