package epp

import (
	"github.com/jtejido/ngac/pkg/pip/graph"
	"github.com/jtejido/ngac/pkg/pip/obligations"
	"github.com/jtejido/ngac/pkg/pip/prohibitions"
)

type FunctionExecutor interface {

	/**
	 * The name of the function
	 * @return the name of the function.
	 */
	Name() string

	/**
	 * How many parameters are expected.
	 * @return the number of parameters this function expects
	 */
	NumParams() int

	/**
	 * Execute the function.
	 * @param graph the graph
	 * @param prohibitions the prohibitions
	 * @param obligations the obligations
	 * @param eventCtx the event that is being processed
	 * @param function the function information
	 * @param functionEvaluator a FunctionEvaluator to evaluate a nested functions
	 * @return the object that the function is expected to return
	 * @throws PMException if there is any error executing the function
	 */
	Exec(graph graph.Graph, prohibitions prohibitions.Prohibitions, obligations obligations.Obligations,
		eventCtx EventContext, function *obligations.Function, functionEvaluator *FunctionEvaluator) (interface{}, error)
}
