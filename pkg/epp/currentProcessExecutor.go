package epp

import (
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

type CurrentProcessExecutor struct{}

func (f *CurrentProcessExecutor) Name() string {
	return "current_process"
}
func (f *CurrentProcessExecutor) NumParams() int {
	return 0
}
func (f *CurrentProcessExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
	eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
	return eventCtx.UserCtx().Process(), nil
}
