package epp

import (
	"fmt"
	"ngac/pkg/pip/graph"
	"ngac/pkg/pip/obligations"
	"ngac/pkg/pip/prohibitions"
)

type FunctionEvaluator struct {
	funExecs map[string]FunctionExecutor
}

func NewFunctionEvaluator() *FunctionEvaluator {
	ans := new(FunctionEvaluator)
	ans.funExecs = make(map[string]FunctionExecutor)

	// add the build in functions
	ans.Add(new(ChildOfAssignExecutor))
	ans.Add(new(CreateNodeExecutor))
	ans.Add(new(CurrentProcessExecutor))
	ans.Add(new(CurrentTargetExecutor))
	ans.Add(new(CurrentUserExecutor))
	ans.Add(new(GetChildrenExecutor))
	ans.Add(new(GetNodeExecutor))
	ans.Add(new(GetNodeNameExecutor))
	ans.Add(new(IsNodeContainedInExecutor))
	ans.Add(new(ParentOfAssignExecutor))
	ans.Add(new(ToPropertiesExecutor))
	return ans
}

func (fe *FunctionEvaluator) Add(executor FunctionExecutor) {
	fe.funExecs[executor.Name()] = executor
}

func (fe *FunctionEvaluator) Remove(executor FunctionExecutor) {
	delete(fe.funExecs, executor.Name())
}

func (fe *FunctionEvaluator) FunctionExecutor(name string) FunctionExecutor {
	var f FunctionExecutor
	var ok bool
	if f, ok = fe.funExecs[name]; !ok {
		panic(fmt.Sprintf("%s is not a recognized function", name))
	}
	return f
}

func (fe *FunctionEvaluator) Eval(graph graph.Graph, prohibitions prohibitions.Prohibitions, obligations obligations.Obligations, eventCtx EventContext, function *obligations.Function) (interface{}, error) {
	functionName := function.Name
	functionExecutor := fe.FunctionExecutor(functionName)
	return functionExecutor.Exec(graph, prohibitions, obligations, eventCtx, function, fe)
}
