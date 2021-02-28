package epp

import (
	"github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/jtejido/ngac/pkg/pip/obligations"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
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
