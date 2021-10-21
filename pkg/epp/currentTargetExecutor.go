package epp

import (
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
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
