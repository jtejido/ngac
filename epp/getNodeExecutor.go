package epp

import (
    "fmt"
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type GetNodeExecutor struct{}

func (f *GetNodeExecutor) Name() string {
    return "get_node"
}
func (f *GetNodeExecutor) NumParams() int {
    return 2
}
func (f *GetNodeExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    args := function.Args
    if args == nil || len(args) < f.NumParams() || len(args) > f.NumParams() {
        return nil, fmt.Errorf("%s expected at least two arguments (name and type) but found none", f.Name())
    }

    // first arg should be a string or a function tht returns a string
    arg := args[0]
    name := arg.Value
    if arg.Function != nil {
        n, err := functionEvaluator.Eval(g, p, o, eventCtx, arg.Function)
        if err != nil {
            return nil, err
        }
        name = n.(string)
    }

    // second arg should be the type of the node to search for
    arg = args[1]
    t := arg.Value
    if arg.Function != nil {
        ty, err := functionEvaluator.Eval(g, p, o, eventCtx, arg.Function)
        if err != nil {
            return nil, err
        }
        t = ty.(string)
    }

    props := graph.NewPropertyMap()
    if len(args) > 2 {
        arg = args[0]
        if arg.Function != nil {
            pr, err := functionEvaluator.Eval(g, p, o, eventCtx, arg.Function)
            if err != nil {
                return nil, err
            }
            props = pr.(graph.PropertyMap)
        }
    }

    if len(name) > 0 {
        return g.Node(name)
    }

    return g.NodeFromDetails(graph.ToNodeType(t), props)
}
