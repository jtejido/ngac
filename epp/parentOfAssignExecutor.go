package epp

import (
    "fmt"
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type ParentOfAssignExecutor struct{}

func (f *ParentOfAssignExecutor) Name() string {
    return "parent_of_assign"
}
func (f *ParentOfAssignExecutor) NumParams() int {
    return 0
}
func (f *ParentOfAssignExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    var parent *graph.Node
    if _, ok := eventCtx.(*AssignToEvent); ok {
        parent = eventCtx.Target()
    } else if v, ok := eventCtx.(*AssignEvent); ok {
        parent = v.ParentNode
    } else if _, ok := eventCtx.(*DeassignFromEvent); ok {
        parent = eventCtx.Target()
    } else if v, ok := eventCtx.(*DeassignEvent); ok {
        parent = v.ParentNode
    } else {
        return nil, fmt.Errorf("invalid event context for function parent_of_assign. Valid event contexts are AssignTo,  Assign, DeassignFrom, and Deassign")
    }

    return parent, nil
}
