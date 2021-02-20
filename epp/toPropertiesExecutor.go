package epp

import (
    "github.com/jtejido/ngac/pip/graph"
    "github.com/jtejido/ngac/pip/obligations"
    "github.com/jtejido/ngac/pip/prohibitions"
    "strings"
)

type ToPropertiesExecutor struct{}

func (f *ToPropertiesExecutor) Name() string {
    return "to_props"
}
func (f *ToPropertiesExecutor) NumParams() int {
    return 0
}
func (f *ToPropertiesExecutor) Exec(g graph.Graph, p prohibitions.Prohibitions, o obligations.Obligations,
    eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error) {
    props := make(map[string]string)
    for _, arg := range function.Args {
        value := arg.Value
        tokens := strings.Split(value, "=")
        if len(tokens) == 2 {
            props[tokens[0]] = tokens[1]
        }
    }

    return props, nil
}
