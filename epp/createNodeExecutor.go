package epp

import (
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
)

type CreateNodeExecutor struct{}

func (f *CreateNodeExecutor) Name() string {
    return "create_node"
}

/**
 * parent name, parent type, parent properties, name, type, properties
 * @return
 */
func (f *CreateNodeExecutor) NumParams() int {
    return 6
}
func (f *CreateNodeExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    args := function.Args

    // first arg is the name, can be function that returns a string
    parentNameArg := args[0]
    parentName := parentNameArg.Value
    if parentNameArg.Function != nil {
        p, err := functionEvaluator.Eval(g, p, o, eventCtx, parentNameArg.Function)
        if err != nil {
            return nil, err
        }
        parentName = p.(string)
    }

    // second arg is the type, can be function
    parentTypeArg := args[1]
    parentType := parentTypeArg.Value
    if parentTypeArg.Function != nil {
        pt, err := functionEvaluator.Eval(g, p, o, eventCtx, parentTypeArg.Function)
        if err != nil {
            return nil, err
        }
        parentType = pt.(string)
    }

    // fourth arg is the name, can be function
    nameArg := args[2]
    name := nameArg.Value
    if nameArg.Function != nil {
        n, err := functionEvaluator.Eval(g, p, o, eventCtx, nameArg.Function)
        if err != nil {
            return nil, err
        }
        name = n.(string)
    }

    // fifth arg is the type, can be function
    typeArg := args[3]
    t := typeArg.Value
    if typeArg.Function != nil {
        tt, err := functionEvaluator.Eval(g, p, o, eventCtx, typeArg.Function)
        if err != nil {
            return nil, err
        }
        t = tt.(string)
    }

    // sixth arg is the properties which is a map that has to come from a function
    props := graph.NewPropertyMap()
    if len(args) > 4 {
        propsArg := args[4]
        if propsArg.Function != nil {
            ppt, err := functionEvaluator.Eval(g, p, o, eventCtx, propsArg.Function)
            if err != nil {
                return nil, err
            }
            props = ppt.(graph.PropertyMap)
        }
    }

    var parentNode *graph.Node
    // if (len(parentName) != 0) {
    parentNode, err := g.Node(parentName)
    if err != nil {
        return nil, err
    }
    // } else {
    //     parentNode = graph.getNode(NodeType.toNodeType(parentType), new HashMap<>());
    // }

    return g.CreateNode(name, graph.ToNodeType(t), props, parentNode.Name)
}
