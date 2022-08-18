package epp

import (
    "github.com/jtejido/ngac/pkg/pip/graph"
    "github.com/jtejido/ngac/pkg/pip/obligations"
    "github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type GetChildrenExecutor struct{}

func (f *GetChildrenExecutor) Name() string {
    return "get_children"
}
func (f *GetChildrenExecutor) NumParams() int {
    return 1
}
func (f *GetChildrenExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    getNodeExecutor := functionEvaluator.FunctionExecutor("get_node")
    node, err := getNodeExecutor.Exec(g, p, o, eventCtx, function, functionEvaluator)
    if err != nil {
        return nil, err
    }

    children := g.Children(node.(*graph.Node).Name)
    return children, nil
}
