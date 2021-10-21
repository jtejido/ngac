package epp

import (
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

type CurrentUserExecutor struct{}

func (f *CurrentUserExecutor) Name() string {
	return "current_user"
}
func (f *CurrentUserExecutor) NumParams() int {
	return 0
}
func (f *CurrentUserExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
	eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
	return g.Node(eventCtx.UserCtx().User())
}
