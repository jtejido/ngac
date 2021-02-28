package epp

import (
    "fmt"
    "github.com/jtejido/ngac/pkg/pip/graph"
    "github.com/jtejido/ngac/pkg/pip/obligations"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type GetNodeNameExecutor struct{}

func (f *GetNodeNameExecutor) Name() string {
    return "get_node_name"
}
func (f *GetNodeNameExecutor) NumParams() int {
    return 1
}
func (f *GetNodeNameExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    args := function.Args
    if len(args) != f.NumParams() {
        return nil, fmt.Errorf("%s expected %d arg but got %d", f.Name(), f.NumParams(), len(args))
    }

    arg := args[0]
    argFunction := arg.Function
    if argFunction == nil {
        return nil, fmt.Errorf("%s expected the first argument to be a function but it was null", f.Name())
    }

    n, err := functionEvaluator.Eval(g, p, o, eventCtx, argFunction)
    if err != nil {
        return nil, err
    }

    node := n.(*graph.Node)
    return node.Name, nil
}
