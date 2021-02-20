package epp

import (
    "fmt"
    "github.com/jtejido/ngac/internal/set"
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type IsNodeContainedInExecutor struct{}

func (f *IsNodeContainedInExecutor) Name() string {
    return "is_node_contained_in"
}
func (f *IsNodeContainedInExecutor) NumParams() int {
    return 2
}
func (f *IsNodeContainedInExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    args := function.Args
    if len(args) != f.NumParams() {
        return nil, fmt.Errorf("%s expected %d arg but got %d", f.Name(), f.NumParams(), len(args))
    }

    arg := args[0]
    ff := arg.Function
    if ff == nil {
        return nil, fmt.Errorf("%s expects two functions as parameters", ff.Name)
    }

    cn, err := functionEvaluator.Eval(g, p, o, eventCtx, ff)
    if err != nil {
        return nil, err
    }
    childNode := cn.(*graph.Node)
    if childNode == nil {
        return false, nil
    }

    arg = args[1]
    ff = arg.Function
    if ff == nil {
        return nil, fmt.Errorf("%s expects two functions as parameters", f.Name())
    }

    pn, err := functionEvaluator.Eval(g, p, o, eventCtx, ff)
    if err != nil {
        return nil, err
    }
    parentNode := pn.(*graph.Node)
    if parentNode == nil {
        return false, nil
    }

    dfs := graph.NewDFS(g)
    nodes := set.NewSet()
    visitor := func(node *graph.Node) error {
        if node.Name == parentNode.Name {
            nodes.Add(node.Name)
        }

        return nil
    }
    propagator := func(parent, child *graph.Node) error { return nil }
    err = dfs.Traverse(childNode, propagator, visitor, graph.PARENTS)
    if err != nil {
        return nil, err
    }
    return nodes.Contains(parentNode.Name), nil
}
