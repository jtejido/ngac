package epp

import (
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type CurrentTargetExecutor struct{}

func (f *CurrentTargetExecutor) Name() string {
    return "current_target"
}
func (f *CurrentTargetExecutor) NumParams() int {
    return 0
}
func (f *CurrentTargetExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    return eventCtx.Target(), nil
}
